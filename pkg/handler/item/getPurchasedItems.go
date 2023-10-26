package item_handler

import (
	"ketalk-api/common"
	item_manager "ketalk-api/pkg/manager/item"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetPurchasedItemsResponse struct {
	Items []ItemBlock `json:"items"`
}

func (h *HttpHandler) GetPurchasedItems(ctx *gin.Context, r *http.Request) (interface{}, error) {
	resp, err := h.handler.GetPurchasedItems(ctx)
	return resp, err
}

func (h *handler) GetPurchasedItems(ctx *gin.Context) (*GetPurchasedItemsResponse, error) {
	userID, err := common.GetUserId(ctx.Request.Context())
	if err != nil {
		return nil, err
	}
	req := item_manager.GetPurchasedItemsRequest{
		UserID: userID,
	}
	purchasedItems, err := h.manager.GetPurchasedItems(ctx, req)
	if err != nil {
		return nil, err
	}

	var items []ItemBlock = make([]ItemBlock, len(purchasedItems))
	for i, item := range purchasedItems {
		items[i] = ItemBlock{
			ID:            item.ID,
			Title:         item.Title,
			Description:   item.Description,
			Price:         item.Price,
			OwnerID:       item.OwnerID,
			FavoriteCount: item.FavoriteCount,
			MessageCount:  item.MessageCount,
			SeenCount:     item.SeenCount,
			ItemStatus:    string(item.ItemStatus),
			CreatedAt:     item.CreatedAt.UTC().Unix(),
			Thumbnail:     item.Thumbnail,
			IsHidden:      item.IsHidden,
		}
	}
	return &GetPurchasedItemsResponse{
		Items: items,
	}, nil
}
