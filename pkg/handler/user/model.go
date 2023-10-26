package user_handler

import (
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	GetUser(ctx *gin.Context) (*GetUserResponse, error)
	UpdateUser(ctx *gin.Context, req UpdateUserRequest) (*UpdateUserResponse, error)
	GetPresignedUrl(ctx *gin.Context) (*GetPresignedUrlResponse, error)
}
