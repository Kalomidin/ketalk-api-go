package item_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DeleteItemResponse struct {
	Success bool `json:"success"`
}

func (h *HttpHandler) DeleteItem(ctx *gin.Context, r *http.Request) (interface{}, error) {
	resp, err := h.handler.DeleteItem(ctx)
	return resp, err
}

func (h *handler) DeleteItem(ctx *gin.Context) (*DeleteItemResponse, error) {
	itemID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return nil, err
	}

	err = h.manager.DeleteItem(ctx, itemID)
	if err != nil {
		return nil, err
	}

	return &DeleteItemResponse{
		Success: true,
	}, nil
}
