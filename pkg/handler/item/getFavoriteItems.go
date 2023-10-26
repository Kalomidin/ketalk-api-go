package item_handler

import (
	"ketalk-api/common"
	item_manager "ketalk-api/pkg/manager/item"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetFavoriteItemsResponse struct {
	Items []ItemBlock `json:"items"`
}

func (h *HttpHandler) GetFavoriteItems(ctx *gin.Context, r *http.Request) (interface{}, error) {
	resp, err := h.handler.GetFavoriteItems(ctx)
	return resp, err
}

func (h *handler) GetFavoriteItems(ctx *gin.Context) (*GetFavoriteItemsResponse, error) {
	userID, err := common.GetUserId(ctx.Request.Context())
	if err != nil {
		return nil, err
	}
	req := item_manager.GetFavoriteItemsRequest{
		UserID: userID,
	}
	resp, err := h.manager.GetFavoriteItems(ctx, req)
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
	return &GetFavoriteItemsResponse{
		Items: items,
	}, nil
}
