package item_handler

import (
	"ketalk-api/common"
	item_manager "ketalk-api/pkg/manager/item"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetSimilarItemsResponse struct {
	SuggestedItems []ItemBlock `json:"suggestedItems"`
	OtherUserItems []ItemBlock `json:"otherUserItems"`
}

func (h *HttpHandler) GetSimilarItems(ctx *gin.Context, r *http.Request) (interface{}, error) {
	resp, err := h.handler.GetSimilarItems(ctx)
	return resp, err
}

func (h *handler) GetSimilarItems(ctx *gin.Context) (*GetSimilarItemsResponse, error) {
	itemID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return nil, err
	}
	userID, err := common.GetUserId(ctx.Request.Context())
	if err != nil {
		return nil, err
	}
	req := item_manager.GetSimilarItemsRequest{
		ItemID: itemID,
		UserID: userID,
	}
	resp, err := h.manager.GetSimilarItems(ctx, req)
	if err != nil {
		return nil, err
	}

	var suggestedItems []ItemBlock = make([]ItemBlock, len(resp.SuggestedItems))
	for i, item := range resp.SuggestedItems {
		suggestedItems[i] = ItemBlock{
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

	var otherUserItems []ItemBlock = make([]ItemBlock, len(resp.OtherUserItems))
	for i, item := range resp.OtherUserItems {
		otherUserItems[i] = ItemBlock{
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

	return &GetSimilarItemsResponse{
		SuggestedItems: suggestedItems,
		OtherUserItems: otherUserItems,
	}, nil
}
