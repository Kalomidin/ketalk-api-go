package item_handler

import (
	"ketalk-api/common"
	item_manager "ketalk-api/pkg/manager/item"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Item struct {
	ID             uuid.UUID `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Price          uint32    `json:"price"`
	Owner          Owner     `json:"owner"`
	ItemStatus     string    `json:"itemStatus"`
	IsHidden       bool      `json:"isHidden"`
	IsUserFavorite bool      `json:"isUserFavorite"`
	Negotiable     bool      `json:"negotiable"`
	FavoriteCount  uint32    `json:"favoriteCount"`
	MessageCount   uint32    `json:"messageCount"`
	SeenCount      uint32    `json:"seenCount"`
	CreatedAt      int64     `json:"createdAt"`
	Thumbnail      string    `json:"thumbnail"`
	Images         []string  `json:"images"`
}

type Owner struct {
	ID       uuid.UUID        `json:"id"`
	Name     string           `json:"name"`
	Avatar   *string          `json:"avatar"`
	Location *common.Location `json:"location"`
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
	return &Item{
		ID:          resp.ID,
		Title:       resp.Title,
		Description: resp.Description,
		Price:       resp.Price,
		Owner: Owner{
			ID:       resp.Owner.ID,
			Name:     resp.Owner.Name,
			Avatar:   resp.Owner.Avatar,
			Location: resp.Owner.Location,
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
		Images:         resp.Images,
	}, nil
}
