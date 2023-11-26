package item_handler

import (
	"ketalk-api/common"
	item_manager "ketalk-api/pkg/manager/item"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetItemsResponse struct {
	Items []ItemBlock `json:"items"`
}

type ItemBlock struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Price         uint32    `json:"price"`
	OwnerID       uuid.UUID `json:"ownerId"`
	ItemStatus    string    `json:"itemStatus"`
	CreatedAt     int64     `json:"createdAt"`
	Thumbnail     string    `json:"thumbnail"`
	IsHidden      bool      `json:"isHidden"`
	FavoriteCount uint32    `json:"favoriteCount"`
	MessageCount  uint32    `json:"messageCount"`
	SeenCount     uint32    `json:"seenCount"`
}

func (h *HttpHandler) GetItems(ctx *gin.Context, r *http.Request) (interface{}, error) {
	resp, err := h.handler.GetItems(ctx)
	return resp, err
}

func (h *handler) GetItems(ctx *gin.Context) (*GetItemsResponse, error) {
	location, err := common.GetLocation(ctx.Request)
	if err != nil {
		return nil, err
	}
	userID, err := common.GetUserId(ctx.Request.Context())
	if err != nil {
		return nil, err
	}

	req := item_manager.GetItemsRequest{
		Location: *location,
		UserID:   userID,
	}
	resp, err := h.manager.GetItems(ctx, req)
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
	return &GetItemsResponse{
		Items: items,
	}, nil
}
