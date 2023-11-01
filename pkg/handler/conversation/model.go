package conversation_handler

import "github.com/gin-gonic/gin"

type ConversationHandler interface {
	CreateConversation(ctx *gin.Context, req CreateConversationRequest) (*CreateConversationResponse, error)
	GetConversations(ctx *gin.Context, req GetConversationsRequest) (*GetConversationsResponse, error)
	GetMembers(ctx *gin.Context) (*GetMembersResponse, error)
}
