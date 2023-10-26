package user_handler

import (
	"ketalk-api/common"
	user_manager "ketalk-api/pkg/manager/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetPresignedUrlResponse struct {
	Url       string `json:"url"`
	ImageName string `json:"imageName"`
}

func (h *HttpHandler) GetPresignedUrl(ctx *gin.Context, r *http.Request) (interface{}, error) {
	resp, err := h.handler.GetPresignedUrl(ctx)
	return resp, err
}

func (h *handler) GetPresignedUrl(ctx *gin.Context) (*GetPresignedUrlResponse, error) {
	userID, err := common.GetUserId(ctx.Request.Context())
	if err != nil {
		return nil, err
	}
	req := user_manager.GetPresignedUrlRequest{
		UserID: userID,
	}
	resp, err := h.manager.GetPresignedUrl(ctx, req)
	if err != nil {
		return nil, err
	}
	return &GetPresignedUrlResponse{
		Url:       resp.Url,
		ImageName: resp.ImageName,
	}, nil
}
