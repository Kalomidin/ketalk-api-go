package ws

import (
	"context"
	"fmt"
	"ketalk-api/common"
	"ketalk-api/common/response"
	"ketalk-api/jwt"
	conn_redis "ketalk-api/pkg/manager/conversation/redis"
	con_repo "ketalk-api/pkg/manager/conversation/repository"
	"ketalk-api/pkg/manager/port"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WebSocketServer struct {
	receiver    *Receiver
	upgrader    *websocket.Upgrader
	router      *gin.Engine
	messageRepo con_repo.MessageRepository
	memberRepo  con_repo.MemberRepository
	userPort    port.UserPort
}

func NewWebSocketServer(
	ctx context.Context,
	userPort port.UserPort,
	messageRepo con_repo.MessageRepository,
	memberRepo con_repo.MemberRepository,
	middleware common.Middleware,
	cfg jwt.Config,
	redis conn_redis.RedisClient,
) (*WebSocketServer, error) {
	r := gin.Default()

	r.Use(middleware.AuthMiddleware(cfg))

	webSocketServer := WebSocketServer{
		receiver: NewReceiver(ctx, redis),
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		router:      r,
		messageRepo: messageRepo,
		memberRepo:  memberRepo,
		userPort:    userPort,
	}

	r.Handle("GET", "/:conversationId", func(ctx *gin.Context) {
		_, err := middleware.HandlerWithAuth(webSocketServer.serveWs)(ctx, ctx.Request)
		if err != nil {
			fmt.Printf("failed to serve web socket connection: %+v\n", err)
			sendErr := response.NewError(err, http.StatusInternalServerError).Send(ctx.Writer)
			if sendErr != nil {
				log.Printf("failed to send failure response for web socket conn: %v, failure: %+v", sendErr, err)
			}
		}
	})
	return &webSocketServer, nil
}

func (ws *WebSocketServer) Serve(port int) error {
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: ws.router,
	}

	fmt.Printf("websocket server is running on port %d\n", port)
	if err := srv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *WebSocketServer) serveWs(ctx *gin.Context, req *http.Request) (interface{}, error) {
	fmt.Printf("serving web socket connection\n")
	conversationId, err := uuid.Parse(ctx.Param("conversationId"))
	if err != nil {
		return nil, err
	}
	userId, err := common.GetUserId(req.Context())
	if err != nil {
		return nil, err
	}
	user, err := s.userPort.GetUser(req.Context(), userId)
	if err != nil {
		return nil, err
	}

	// check if user is member of conversation
	if members, err := s.memberRepo.GetMembers(ctx, conversationId); err != nil {
		return nil, err
	} else {
		isMember := false
		for _, member := range members {
			if member.MemberID == userId {
				isMember = true
				break
			}
		}
		if !isMember {
			return nil, fmt.Errorf("user is not a member of conversation")
		}
	}

	ws, err := s.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Printf("failed while upgrade: %+v\n", err)
		return nil, err
	}

	client := &Client{
		conn:           ws,
		UserID:         userId,
		ConversationID: conversationId,
		UserName:       user.Username,
		GroupID:        s.receiver.redis.GetGroupID(conversationId),
	}

	// get the messages
	mes, err := s.messageRepo.GetMessages(ctx, conversationId)
	if err != nil {
		return nil, err
	}
	var messages []ServerToActorMessage = make([]ServerToActorMessage, len(mes))
	for i, m := range mes {
		messages[i] = ServerToActorMessage{
			Message:     m.Message,
			SenderID:    m.SenderID,
			CreatedAt:   m.CreatedAt.UTC().Unix(),
			MessageType: MessageTypeMessage,
		}
	}

	if err := client.conn.WriteJSON(ServerToActorMessages{
		Messages: messages,
	}); err != nil {
		fmt.Printf("failed to return old messages to client: %v\n", err)
	}

	if err = s.receiver.Add(userId, conversationId, client); err != nil {
		return nil, err
	}
	if err = s.receiver.Receive(ctx, client); err != nil {
		fmt.Printf("client receiver failed while receiving: %+v\n", err)
	}

	s.receiver.Remove(userId, conversationId)

	return nil, nil
}
