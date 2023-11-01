package conversation_handler

import (
	"ketalk-api/common"
	conversation_manager "ketalk-api/pkg/manager/conversation"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetConversationsRequest struct {
}

type GetConversationsResponse struct {
	Conversations []Conversation `json:"conversations"`
}

type Conversation struct {
	Title                 string `json:"title"`
	LastMessage           string `json:"lastMessage"`
	LastMessageAt         int64  `json:"lastMessageAt"`
	ID                    string `json:"id"`
	LastMessageSenderID   string `json:"lastMessageSenderId"`
	SecondaryUserImageUrl string `json:"secondaryUserImageUrl"`
	ItemID                string `json:"itemId"`
	ItemThumbnail         string `json:"itemThumbnail"`
	IsMessageRead         bool   `json:"isMessageRead"`
}

func (h *HttpHandler) GetConversations(ctx *gin.Context, req *http.Request) (interface{}, error) {
	var request GetConversationsRequest
	return h.handler.GetConversations(ctx, request)
}

func (h *handler) GetConversations(ctx *gin.Context, req GetConversationsRequest) (*GetConversationsResponse, error) {
	userID, err := common.GetUserId(ctx.Request.Context())
	if err != nil {
		return nil, err
	}
	var request = conversation_manager.GetConversationsRequest{
		UserID: userID,
	}
	resp, err := h.manager.GetConversations(ctx, request)
	if err != nil {
		return nil, err
	}
	var conversations []Conversation = make([]Conversation, len(resp))
	for i, conversation := range resp {
		conversations[i] = Conversation{
			Title:                 conversation.Title,
			LastMessage:           conversation.LastMessage,
			LastMessageAt:         conversation.LastMessageAt.UTC().Unix(),
			ID:                    conversation.ID.String(),
			LastMessageSenderID:   conversation.LastMessageSenderID.String(),
			SecondaryUserImageUrl: conversation.SecondaryUserImageUrl,
			ItemID:                conversation.ItemID.String(),
			ItemThumbnail:         conversation.ItemThumbnail,
			IsMessageRead:         conversation.IsMessageRead,
		}
	}
	return &GetConversationsResponse{
		Conversations: conversations,
	}, nil
}
