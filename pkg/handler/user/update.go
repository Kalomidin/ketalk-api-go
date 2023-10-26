package user_handler

import (
	"ketalk-api/common"
	user_manager "ketalk-api/pkg/manager/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateUserRequest struct {
	Name  *string `json:"name"`
	Image *string `json:"image"`
}

type UpdateUserResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Image    *string   `json:"image"`
}

func (h *HttpHandler) UpdateUser(ctx *gin.Context, r *http.Request) (interface{}, error) {
	var req UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	resp, err := h.handler.UpdateUser(ctx, req)
	return resp, err
}

func (h *handler) UpdateUser(ctx *gin.Context, req UpdateUserRequest) (*UpdateUserResponse, error) {
	userID, err := common.GetUserId(ctx.Request.Context())
	if err != nil {
		return nil, err
	}
	user, err := h.manager.Update(ctx, user_manager.UpdateUserRequest{
		UserID: userID,
		Name:   req.Name,
		Image:  req.Image,
	})
	if err != nil {
		return nil, err
	}
	return &UpdateUserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Image:    user.Image,
	}, nil
}
