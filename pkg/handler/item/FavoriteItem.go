package item_handler

import (
	"ketalk-api/common"
	item_manager "ketalk-api/pkg/manager/item"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FavoriteItemRequest struct {
	IsFavorite bool `json:"isFavorite"`
}

type FavoriteItemResponse struct {
}

func (h *HttpHandler) FavoriteItem(ctx *gin.Context, r *http.Request) (interface{}, error) {
	var req FavoriteItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	resp, err := h.handler.FavoriteItem(ctx, req)
	return resp, err
}

func (h *handler) FavoriteItem(ctx *gin.Context, req FavoriteItemRequest) (*FavoriteItemResponse, error) {
	userID, err := common.GetUserId(ctx.Request.Context())
	if err != nil {
		return nil, err
	}
	itemID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return nil, err
	}
	_req := item_manager.FavoriteItemRequest{
		UserID:     userID,
		ItemID:     itemID,
		IsFavorite: req.IsFavorite,
	}
	if _, err := h.manager.FavoriteItem(ctx, _req); err != nil {
		return nil, err
	}
	return &FavoriteItemResponse{}, nil
}
