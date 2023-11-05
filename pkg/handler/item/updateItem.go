package item_handler

import (
	"ketalk-api/common"
	item_manager "ketalk-api/pkg/manager/item"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateItemRequest struct {
	IsHidden    *bool              `json:"isHidden"`
	ItemStatus  *string            `json:"itemStatus"`
	Title       *string            `json:"title"`
	Description *string            `json:"description"`
	Price       *uint32            `json:"price"`
	Negotiable  *bool              `json:"negotiable"`
	Size        *float32           `json:"size"`
	Weight      *float32           `json:"weight"`
	KaratId     *uuid.UUID         `json:"karatId"`
	CategoryId  *uuid.UUID         `json:"categoryId"`
	Images      []UpdatedItemImage `json:"images"`
}

type UpdatedItemImage struct {
	ID      *uuid.UUID `json:"id"`
	Name    *string    `json:"name"`
	IsCover bool       `json:"isCover"`
}

type UpdateItemResponse struct {
	NewImagesPresignedUrls []ImageUploadUrlWithName `json:"newImagesPresignedUrls"`
}

func (h *HttpHandler) UpdateItem(ctx *gin.Context, r *http.Request) (interface{}, error) {
	var req UpdateItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	resp, err := h.handler.UpdateItem(ctx, req)
	return resp, err
}

func (h *handler) UpdateItem(ctx *gin.Context, req UpdateItemRequest) (*UpdateItemResponse, error) {
	userID, err := common.GetUserId(ctx.Request.Context())
	if err != nil {
		return nil, err
	}
	itemID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return nil, err
	}

	var itemStatus *item_manager.ItemStatus
	if req.ItemStatus != nil {
		itemStatus, err = item_manager.ParseItemStatus(*req.ItemStatus)
		if err != nil {
			return nil, err
		}
	}
	var images []item_manager.UpdatedItemImage = make([]item_manager.UpdatedItemImage, len(req.Images))
	for i, image := range req.Images {
		images[i] = item_manager.UpdatedItemImage{
			ID:      image.ID,
			Name:    image.Name,
			IsCover: image.IsCover,
		}
	}
	updateItemReq := item_manager.UpdateItemRequest{
		UserID:      userID,
		ItemID:      itemID,
		IsHidden:    req.IsHidden,
		ItemStatus:  itemStatus,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Negotiable:  req.Negotiable,
		Size:        req.Size,
		Weight:      req.Weight,
		KaratID:     req.KaratId,
		CategoryID:  req.CategoryId,
		Images:      images,
	}
	resp, err := h.manager.UpdateItem(ctx, updateItemReq)
	if err != nil {
		return nil, err
	}
	var newImagesPresignedUrls []ImageUploadUrlWithName = make([]ImageUploadUrlWithName, len(resp.NewImagesPresignedUrls))

	for i, image := range resp.NewImagesPresignedUrls {
		newImagesPresignedUrls[i] = ImageUploadUrlWithName{
			ID:        image.ID,
			Name:      image.Name,
			SignedUrl: image.SignedUrl,
		}
	}
	return &UpdateItemResponse{
		NewImagesPresignedUrls: newImagesPresignedUrls,
	}, nil
}
