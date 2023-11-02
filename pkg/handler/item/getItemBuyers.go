package item_handler

import (
	item_manager "ketalk-api/pkg/manager/item"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetItemBuyersResponse struct {
	Buyers []ItemBuyer `json:"buyers"`
}

type ItemBuyer struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Avatar         *string   `json:"avatar"`
	LastMessagedAt int64     `json:"lastMessagedAt"`
}

func (h *HttpHandler) GetItemBuyers(ctx *gin.Context, r *http.Request) (interface{}, error) {
	resp, err := h.handler.GetItemBuyers(ctx)
	return resp, err
}

func (h *handler) GetItemBuyers(ctx *gin.Context) (*GetItemBuyersResponse, error) {
	itemId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return nil, err
	}
	req := item_manager.GetItemBuyersRequest{
		ItemID: itemId,
	}
	resp, err := h.manager.GetItemBuyers(ctx, req)
	if err != nil {
		return nil, err
	}
	var buyers []ItemBuyer = make([]ItemBuyer, len(resp))
	for i, buyer := range resp {
		buyers[i] = ItemBuyer{
			ID:             buyer.ID,
			Name:           buyer.Name,
			Avatar:         buyer.Avatar,
			LastMessagedAt: buyer.LastMessagedAt.Unix(),
		}
	}

	return &GetItemBuyersResponse{
		Buyers: buyers,
	}, nil
}
