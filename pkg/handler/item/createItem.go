package item_handler

import (
	"ketalk-api/common"
	item_manager "ketalk-api/pkg/manager/item"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateItemRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Negotiable  bool      `json:"negotiable"`
	Price       uint32    `json:"price"`
	Size        float32   `json:"size"`
	Weight      float32   `json:"weight"`
	KaratID     uuid.UUID `json:"karatId"`
	CategoryID  uuid.UUID `json:"categoryId"`
	Images      []string  `json:"images"`
	Thumbnail   string    `json:"thumbnail"`
}

type ItemImage struct {
	ID        uuid.UUID `json:"id"`
	SignedUrl string    `json:"url"`
	Name      string    `json:"name"`
}

type ImageUploadUrlWithName struct {
	ID        uuid.UUID `json:"id"`
	SignedUrl string    `json:"url"`
	Name      string    `json:"name"`
}

type CreateItemResponse struct {
	ID            uuid.UUID                `json:"id"`
	CreatedAt     int64                    `json:"createdAt"`
	PresignedUrls []ImageUploadUrlWithName `json:"itemImages"`
}

func (h *HttpHandler) CreateItem(ctx *gin.Context, r *http.Request) (interface{}, error) {
	var req CreateItemRequest
	if err := ctx.BindJSON(&req); err != nil {
		return nil, err
	}
	resp, err := h.handler.CreateItem(ctx, req)
	return resp, err
}

func (h *handler) CreateItem(ctx *gin.Context, r CreateItemRequest) (*CreateItemResponse, error) {
	userID, err := common.GetUserId(ctx.Request.Context())
	if err != nil {
		return nil, err
	}
	location, err := common.GetLocation(ctx.Request)
	if err != nil {
		return nil, err
	}

	req := item_manager.AddItemRequest{
		Title:       r.Title,
		Description: r.Description,
		Negotiable:  r.Negotiable,
		Price:       r.Price,
		Size:        r.Size,
		Weight:      r.Weight,
		OwnerID:     userID,
		KaratID:     r.KaratID,
		CategoryID:  r.CategoryID,
		Images:      r.Images,
		Thumbnail:   r.Thumbnail,
		Location:    *location,
	}
	resp, err := h.manager.AddItem(ctx, req)
	if err != nil {
		return nil, err
	}
	var presignedUrls []ImageUploadUrlWithName = make([]ImageUploadUrlWithName, len(resp.PresignedUrls))
	for i, url := range resp.PresignedUrls {
		presignedUrls[i] = ImageUploadUrlWithName{
			ID:        url.ID,
			SignedUrl: url.SignedUrl,
			Name:      url.Name,
		}
	}
	return &CreateItemResponse{
		ID:            resp.ID,
		CreatedAt:     resp.CreatedAt.UTC().Unix(),
		PresignedUrls: presignedUrls,
	}, nil
}
