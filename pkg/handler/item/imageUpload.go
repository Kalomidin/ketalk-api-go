package item_handler

import (
	item_manager "ketalk-api/pkg/manager/item"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadItemImagesRequest struct {
	ItemID   uuid.UUID   `json:"id"`
	ImageIds []uuid.UUID `json:"imageIds"`
}

type UploadItemImagesResponse struct {
}

func (h *HttpHandler) UploadItemImages(ctx *gin.Context, r *http.Request) (interface{}, error) {
	var req UploadItemImagesRequest
	if err := ctx.BindJSON(&req); err != nil {
		return nil, err
	}
	resp, err := h.handler.UploadItemImages(ctx, req)
	return resp, err
}

func (h *handler) UploadItemImages(ctx *gin.Context, r UploadItemImagesRequest) (*UploadItemImagesResponse, error) {
	req := item_manager.UploadItemImagesRequest{
		ItemID:   r.ItemID,
		ImageIds: r.ImageIds,
	}
	if _, err := h.manager.UploadItemImages(ctx, req); err != nil {
		return nil, err
	}
	return &UploadItemImagesResponse{}, nil
}
