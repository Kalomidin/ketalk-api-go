package auth_handler

import (
	"ketalk-api/common"
	auth_manager "ketalk-api/pkg/manager/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func (h *HttpHandler) Logout(ctx *gin.Context, r *http.Request) (interface{}, error) {
	var req LogoutRequest
	if err := ctx.BindJSON(&req); err != nil {
		return nil, err
	}
	if err := h.handler.Logout(ctx, req); err != nil {
		return nil, err
	}
	return nil, nil
}

func (h *handler) Logout(ctx *gin.Context, logoutReq LogoutRequest) error {
	userId, err := common.GetUserId(ctx.Request.Context())
	if err != nil {
		return err
	}

	return h.service.Logout(ctx, auth_manager.LogoutRequest{
		RefreshToken: logoutReq.RefreshToken,
		UserID:       userId,
	})
}
