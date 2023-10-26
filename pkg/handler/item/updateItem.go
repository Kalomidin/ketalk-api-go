package item_handler

import (
	"ketalk-api/common"
	item_manager "ketalk-api/pkg/manager/item"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateItemRequest struct {
	IsHidden    *bool   `json:"isHidden"`
	ItemStatus  *string `json:"itemStatus"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Price       *uint32 `json:"price"`
	Negotiable  *bool   `json:"negotiable"`
}

type UpdateItemResponse struct {
}

func (h *HttpHandler) UpdateItem(ctx *gin.Context, r *http.Request) (interface{}, error) {
	var req UpdateItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	resp, err := h.handler.UpdateItem(ctx, req)
	return resp, err
}

func (h *handler) UpdateItem(ctx *gin.Context, req UpdateItemRequest) (*UpdateItemResponse, error) {
	userID, err := common.GetUserId(ctx.Request.Context())
	if err != nil {
		return nil, err
	}
	itemID, err := uuid.Parse(ctx.Param("id"))

	var itemStatus *item_manager.ItemStatus
	if req.ItemStatus != nil {
		itemStatus, err = item_manager.ParseItemStatus(*req.ItemStatus)
		if err != nil {
			return nil, err
		}
	}
	updateItemReq := item_manager.UpdateItemRequest{
		UserID:      userID,
		ItemID:      itemID,
		IsHidden:    req.IsHidden,
		ItemStatus:  itemStatus,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Negotiable:  req.Negotiable,
	}
	if _, err := h.manager.UpdateItem(ctx, updateItemReq); err != nil {
		return nil, err
	}
	return &UpdateItemResponse{}, nil
}
