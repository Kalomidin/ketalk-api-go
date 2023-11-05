package item_handler

import (
	"ketalk-api/common"
	item_manager "ketalk-api/pkg/manager/item"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Item struct {
	ID             uuid.UUID   `json:"id"`
	Title          string      `json:"title"`
	Description    string      `json:"description"`
	Price          uint32      `json:"price"`
	Owner          Owner       `json:"owner"`
	ItemStatus     string      `json:"itemStatus"`
	IsHidden       bool        `json:"isHidden"`
	IsUserFavorite bool        `json:"isUserFavorite"`
	Negotiable     bool        `json:"negotiable"`
	FavoriteCount  uint32      `json:"favoriteCount"`
	MessageCount   uint32      `json:"messageCount"`
	SeenCount      uint32      `json:"seenCount"`
	CreatedAt      int64       `json:"createdAt"`
	Thumbnail      string      `json:"thumbnail"`
	Images         []ItemImage `json:"images"`
	KaratID        uuid.UUID   `json:"karatId"`
	CategoryID     uuid.UUID   `json:"categoryId"`
	Weigt          float32     `json:"weight"`
	Size           float32     `json:"size"`
}

type Owner struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Avatar   *string   `json:"avatar"`
	Geofence Geofence  `json:"geofence"`
}

type Geofence struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (h *HttpHandler) GetItem(ctx *gin.Context, req *http.Request) (interface{}, error) {
	resp, err := h.handler.GetItem(ctx)
	return resp, err
}

func (h *handler) GetItem(ctx *gin.Context) (*Item, error) {
	itemId, err := uuid.Parse(ctx.Param("id"))

	if err != nil {
		return nil, err
	}
	userID, err := common.GetUserId(ctx.Request.Context())
	if err != nil {
		return nil, err
	}
	req := item_manager.GetItemRequest{
		ItemID: itemId,
		UserID: userID,
	}
	resp, err := h.manager.GetItem(ctx, req)
	if err != nil {
		return nil, err
	}

	var itemImages []ItemImage = make([]ItemImage, len(resp.Images))
	for i, image := range resp.Images {
		itemImages[i] = ItemImage{
			ID:        image.ID,
			SignedUrl: image.SignedUrl,
			Name:      image.Name,
		}
	}

	return &Item{
		ID:          resp.ID,
		Title:       resp.Title,
		Description: resp.Description,
		Price:       resp.Price,
		Owner: Owner{
			ID:     resp.Owner.ID,
			Name:   resp.Owner.Name,
			Avatar: resp.Owner.Avatar,
			Geofence: Geofence{
				ID:   resp.Owner.Geofence.ID,
				Name: resp.Owner.Geofence.Name,
			},
		},
		IsHidden:       resp.IsHidden,
		IsUserFavorite: resp.IsUserFavorite,
		Negotiable:     resp.Negotiable,
		FavoriteCount:  resp.FavoriteCount,
		MessageCount:   resp.MessageCount,
		SeenCount:      resp.SeenCount,
		ItemStatus:     string(resp.ItemStatus),
		CreatedAt:      resp.CreatedAt.Unix(),
		Thumbnail:      resp.Thumbnail,
		Images:         itemImages,
		KaratID:        resp.KaratID,
		CategoryID:     resp.CategoryID,
		Weigt:          resp.Weight,
		Size:           resp.Size,
	}, nil
}
