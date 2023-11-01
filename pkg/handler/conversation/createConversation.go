package conversation_handler

import (
	"ketalk-api/common"
	conversation_manager "ketalk-api/pkg/manager/conversation"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateConversationRequest struct {
	ItemID uuid.UUID `json:"itemId"`
}

type CreateConversationResponse struct {
	ID              uuid.UUID `json:"id"`
	SecondaryUserID uuid.UUID `json:"secondaryUserId"`
	ItemID          uuid.UUID `json:"itemId"`
}

func (h *HttpHandler) CreateConversation(ctx *gin.Context, req *http.Request) (interface{}, error) {
	var request CreateConversationRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		return nil, err
	}
	return h.handler.CreateConversation(ctx, request)
}

func (h *handler) CreateConversation(ctx *gin.Context, req CreateConversationRequest) (*CreateConversationResponse, error) {
	userID, err := common.GetUserId(ctx.Request.Context())
	if err != nil {
		return nil, err
	}
	var request = conversation_manager.CreateConversationRequest{
		ItemID: req.ItemID,
		UserID: userID,
	}
	resp, err := h.manager.CreateConversation(ctx, request)
	if err != nil {
		return nil, err
	}
	return &CreateConversationResponse{
		ID:              resp.ID,
		SecondaryUserID: resp.SecondaryUserID,
		ItemID:          resp.ItemID,
	}, nil
}
