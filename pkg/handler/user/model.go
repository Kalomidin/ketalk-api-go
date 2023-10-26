package user_handler

import (
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	GetUser(ctx *gin.Context) (*GetUserResponse, error)
}
