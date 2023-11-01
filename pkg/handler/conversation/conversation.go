package conversation_handler

import (
	"context"
	"fmt"
	"ketalk-api/common"
	conversation_manager "ketalk-api/pkg/manager/conversation"

	"github.com/gin-gonic/gin"
)

type HttpHandler struct {
	handler    ConversationHandler
	middleware common.Middleware
}

type handler struct {
	manager conversation_manager.ConversationManager
}

func NewHttpHandler(handler ConversationHandler, middleware common.Middleware) *HttpHandler {
	return &HttpHandler{
		handler:    handler,
		middleware: middleware,
	}
}

func NewHandler(manager conversation_manager.ConversationManager) ConversationHandler {
	return &handler{
		manager: manager,
	}
}

func (c *HttpHandler) Init(ctx context.Context, router *gin.Engine) {
	routes := map[string]map[string]common.HandlerFunc{
		"POST": {
			"": c.middleware.HandlerWithAuth(c.CreateConversation),
		},
		"GET": {
			"":                         c.middleware.HandlerWithAuth(c.GetConversations),
			"/:conversationId/members": c.middleware.HandlerWithAuth(c.GetMembers),
		},
	}
	for method, route := range routes {
		for r, h := range route {
			router.Handle(method, fmt.Sprintf("/conversation%s", r), common.GenericHandler(h))
		}
	}
	fmt.Println("initialized item handler")
}
