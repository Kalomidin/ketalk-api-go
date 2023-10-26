package auth_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RefreshAccessTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshAccessTokenResponse struct {
	AuthToken    string `json:"authToken"`
	RefreshToken string `json:"refreshToken"`
}

func (h *HttpHandler) RefreshAccessToken(ctx *gin.Context, r *http.Request) (interface{}, error) {
	var req RefreshAccessTokenRequest
	if err := ctx.BindJSON(&req); err != nil {
		return nil, err
	}
	resp, err := h.handler.RefreshAccessToken(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (h *handler) RefreshAccessToken(ctx *gin.Context, req RefreshAccessTokenRequest) (*RefreshAccessTokenResponse, error) {
	return nil, nil
}
