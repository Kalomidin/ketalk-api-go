package item_handler

import (
	"ketalk-api/common"
	item_manager "ketalk-api/pkg/manager/item"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *HttpHandler) IncrementConversationCount(ctx *gin.Context, r *http.Request) (interface{}, error) {
	resp, err := h.handler.IncrementConversationCount(ctx)
	return resp, err
}

func (h *handler) IncrementConversationCount(ctx *gin.Context) (interface{}, error) {
	itemID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return nil, err
	}
	userID, err := common.GetUserId(ctx.Request.Context())
	if err != nil {
		return nil, err
	}
	req := item_manager.IncrementConversationCountRequest{
		ItemID: itemID,
		UserID: userID,
	}
	err = h.manager.IncrementConversationCount(ctx, req)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
