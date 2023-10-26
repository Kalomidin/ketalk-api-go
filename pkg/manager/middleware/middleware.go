package middleware

import (
	"fmt"
	"ketalk-api/common"
	"ketalk-api/jwt"
	"ketalk-api/pkg/manager/port"
	"net/http"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type middleware struct {
	userPort port.UserPort
}

func NewMiddleware(userPort port.UserPort) common.Middleware {
	return &middleware{
		userPort: userPort,
	}
}

func (m *middleware) AuthMiddleware(cfg jwt.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := ctx.Request
		token := req.Header.Get("Authorization")
		if token == "" {
			// continue if no authorization provided
			fmt.Println("continue without token")
			ctx.Next()
			return
		}

		bearer := "Bearer "
		index := strings.Index(token, bearer)
		if index == 0 {
			token = token[len(bearer):]
		}

		contextWithToken, err := jwt.ParseJWTToken(req.Context(), cfg, token)
		if err != nil {
			// todo: we may want to return err if jwt token validation failed
			ctx.Next()
			return
		}
		req = req.WithContext(contextWithToken)
		ctx.Request = req
		ctx.Next()
	}
}

func (m *middleware) ValidateUserAuthorization(ctx *gin.Context) error {
	tokenData, err := jwt.GetJWTToken(ctx.Request.Context())
	if err != nil {
		return err
	}
	userId, err := uuid.Parse(tokenData.Subject)
	if err != nil {
		return fmt.Errorf("invalid token")
	}
	_, err = m.userPort.GetUser(ctx, userId)
	if err != nil {
		return err
	}
	return nil
}

func (m *middleware) HandlerWithAuth(handler common.HandlerFunc) common.HandlerFunc {
	return func(ctx *gin.Context, req *http.Request) (interface{}, error) {
		err := m.ValidateUserAuthorization(ctx)
		if err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}
