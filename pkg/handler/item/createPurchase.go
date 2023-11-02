package item_handler

import (
	item_manager "ketalk-api/pkg/manager/item"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreatePurchaseRequest struct {
	BuyerID uuid.UUID `json:"buyerId"`
}

type CreatePurchaseResponse struct {
}

func (h *HttpHandler) CreatePurchase(ctx *gin.Context, r *http.Request) (interface{}, error) {
	var req CreatePurchaseRequest
	if err := ctx.BindJSON(&req); err != nil {
		return nil, err
	}
	resp, err := h.handler.CreatePurchase(ctx, req)
	return resp, err
}

func (h *handler) CreatePurchase(ctx *gin.Context, req CreatePurchaseRequest) (*CreatePurchaseResponse, error) {
	itemID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return nil, err
	}
	purchaseReq := item_manager.CreatePurchaseRequest{
		ItemID:  itemID,
		BuyerID: req.BuyerID,
	}

	_, err = h.manager.CreatePurchase(ctx, purchaseReq)
	if err != nil {
		return nil, err
	}
	return &CreatePurchaseResponse{}, nil
}
