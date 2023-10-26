package user_handler

import (
	"ketalk-api/common"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetUserResponse struct {
	UserId   uuid.UUID `json:"id"`
	Username string    `json:"userName"`
	Email    string    `json:"email"`
	Image    *string   `json:"avatar"`
}

// @BasePath /api/v1

// GetUser
// @Summary Get user
// @Schemes
// @Description get user
// @Accept json
// @Produce json
// @Param Authorization header string true "Authoriztion"
// @Success 200 {object} GetUserResponse
// @Router /user [get]
func (h *HttpHandler) GetUser(ctx *gin.Context, r *http.Request) (interface{}, error) {
	resp, err := h.handler.GetUser(ctx)
	return resp, err
}

func (h *handler) GetUser(ctx *gin.Context) (*GetUserResponse, error) {
	userID, err := common.GetUserId(ctx.Request.Context())
	if err != nil {
		return nil, err
	}
	user, err := h.service.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &GetUserResponse{
		UserId:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Image:    user.Image,
	}, nil
}
