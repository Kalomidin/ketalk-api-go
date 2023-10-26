package item_handler

import (
	"ketalk-api/common"
	item_manager "ketalk-api/pkg/manager/item"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetUserItemsResponse struct {
	Items []ItemBlock `json:"items"`
}

func (h *HttpHandler) GetUserItems(ctx *gin.Context, r *http.Request) (interface{}, error) {
	resp, err := h.handler.GetUserItems(ctx)
	return resp, err
}

func (h *handler) GetUserItems(ctx *gin.Context) (*GetUserItemsResponse, error) {
	userID, err := common.GetUserId(ctx.Request.Context())
	if err != nil {
		return nil, err
	}
	req := item_manager.GetUserItemsRequest{
		UserID: userID,
	}
	resp, err := h.manager.GetUserItems(ctx, req)
	if err != nil {
		return nil, err
	}
	var items []ItemBlock = make([]ItemBlock, len(resp))
	for i, item := range resp {
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
	return &GetUserItemsResponse{
		Items: items,
	}, nil
}
