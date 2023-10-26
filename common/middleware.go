package common

import (
	"ketalk-api/jwt"

	"github.com/gin-gonic/gin"
)

type Middleware interface {
	AuthMiddleware(cfg jwt.Config) gin.HandlerFunc
	ValidateUserAuthorization(ctx *gin.Context) error
	HandlerWithAuth(handler HandlerFunc) HandlerFunc
}
