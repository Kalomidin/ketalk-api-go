package common

import (
	"context"
	"ketalk-api/jwt"

	"github.com/gin-gonic/gin"
)

type Middleware interface {
	AuthMiddleware(cfg jwt.Config) gin.HandlerFunc
	ValidateUserAuthorization(ctx context.Context) error
	HandlerWithAuth(handler HandlerFunc) HandlerFunc
}
