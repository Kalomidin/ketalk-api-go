package ws

import (
	"context"
	"encoding/json"
	"fmt"
	conn_redis "ketalk-api/pkg/manager/conversation/redis"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	conn           *websocket.Conn
	ConversationID uuid.UUID
	UserID         uuid.UUID
	UserName       string
	GroupID        int
}

type Receiver struct {
	clients map[uuid.UUID]map[uuid.UUID]*Client
	lock    sync.RWMutex
	redis   conn_redis.RedisClient
}

func NewReceiver(ctx context.Context, redis conn_redis.RedisClient) *Receiver {
	receiver := Receiver{
		clients: make(map[uuid.UUID]map[uuid.UUID]*Client),
		redis:   redis,
	}
	// start redis message handler
	go receiver.redis.Handle(ctx, receiver.HandleMessage, "websocket")

	return &receiver
}

func (r *Receiver) Receive(ctx context.Context, client *Client) error {
	for {
		var mes ClientActorMessage
		err := client.conn.ReadJSON(&mes)
		if err != nil {
			// TODO: if invalid format message, maybe continue?
			return err
		}

		log.Printf("received message: %+v\n", mes)
		var msg Message

		switch msg.Type {
		case MessageTypeLeave:
			// TODO: when user leaves, inform other users
			return nil
		default:
			msg.Timestamp = time.Now().UTC().Unix()
			msg.UserID = client.UserID
			msg.ConversationID = client.ConversationID
			msg.Type = mes.Type
			msg.Message = mes.Message
		}

		if err := r.redis.AddMessage(ctx, client.GroupID, client.ConversationID, msg); err != nil {
			fmt.Printf("failed to add message to redis: %v\n", err)
		}
	}
}

func (r *Receiver) HandleMessage(ctx context.Context, payload string) error {
	fmt.Printf("received message from redis: %s\n", payload)
	payloadBytes := []byte(payload)

	var mes Message
	if err := json.Unmarshal(payloadBytes, &mes); err != nil {
		return err
	}

	r.lock.RLock()
	if _, ok := r.clients[mes.ConversationID]; !ok {
		fmt.Printf("no clients for conversationId: %s\n", mes.ConversationID)
		r.lock.RUnlock()
		return fmt.Errorf("no clients for conversationId: %s", mes.ConversationID)
	}
	for _, client := range r.clients[mes.ConversationID] {
		fmt.Printf("sending message to client: %+v\n", client.UserID)
		messages := ServerToActorMessages{
			Messages: []ServerToActorMessage{
				{
					SenderID:    mes.UserID,
					Message:     mes.Message,
					CreatedAt:   mes.Timestamp,
					MessageType: mes.Type,
				},
			},
		}
		fmt.Printf("sending message to client: %+v\n", messages)
		if err := client.conn.WriteJSON(messages); err != nil {
			fmt.Printf("failed to write message to client: %v\n", err)
		}
	}
	r.lock.RUnlock()
	return nil
}

func (r *Receiver) GetClient(userId, conversationId uuid.UUID) *Client {
	r.lock.RLock()
	defer r.lock.RUnlock()
	if _, ok := r.clients[conversationId]; !ok {
		return nil
	}
	if _, ok := r.clients[conversationId][userId]; !ok {
		return nil
	}
	return r.clients[conversationId][userId]
}

func (r *Receiver) Add(userId, conversationId uuid.UUID, client *Client) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	if _, ok := r.clients[conversationId]; !ok {
		r.clients[conversationId] = make(map[uuid.UUID]*Client)
	}
	if _, ok := r.clients[conversationId][userId]; ok {
		return fmt.Errorf("client already exists")
	}
	r.clients[conversationId][userId] = client
	return nil
}

func (r *Receiver) Remove(userId, conversationId uuid.UUID) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	if _, ok := r.clients[conversationId]; !ok {
		return nil
	}
	if _, ok := r.clients[conversationId][userId]; !ok {
		return nil
	}
	delete(r.clients[conversationId], userId)
	if len(r.clients[conversationId]) == 0 {
		delete(r.clients, conversationId)
	}

	return nil
}
