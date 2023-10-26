package item_handler

import (
	"fmt"
	item_manager "ketalk-api/pkg/manager/item"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetItemsResponse struct {
	Items []Item `json:"items"`
}

type Item struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Price         uint32    `json:"price"`
	OwnerID       uuid.UUID `json:"ownerId"`
	FavoriteCount uint32    `json:"favoriteCount"`
	MessageCount  uint32    `json:"messageCount"`
	SeenCount     uint32    `json:"seenCount"`
	ItemStatus    string    `json:"itemStatus"`
	CreatedAt     int64     `json:"createdAt"`
	Thumbnail     string    `json:"thumbnail"`
}

func (h *HttpHandler) GetItems(ctx *gin.Context, r *http.Request) (interface{}, error) {
	resp, err := h.handler.GetItems(ctx)
	return resp, err
}

func (h *handler) GetItems(ctx *gin.Context) (*GetItemsResponse, error) {
	geofenceID, err := strconv.Atoi(ctx.Param("geofenceId"))
	if err != nil {
		return nil, err
	}
	if geofenceID < 0 {
		return nil, fmt.Errorf("invalid geofence id: %d", geofenceID)
	}
	req := item_manager.GetItemsRequest{
		GeofenceID: uint32(geofenceID),
	}
	resp, err := h.manager.GetItems(ctx, req)
	if err != nil {
		return nil, err
	}
	var items []Item = make([]Item, len(resp))
	for i, item := range resp {
		items[i] = Item{
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
		}
	}
	return &GetItemsResponse{
		Items: items,
	}, nil
}
