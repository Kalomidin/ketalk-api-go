package middleware

import (
	"context"
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

func (m *middleware) ValidateUserAuthorization(ctx context.Context) error {
	tokenData, err := jwt.GetJWTToken(ctx)
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
		err := m.ValidateUserAuthorization(req.Context())
		if err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

func (m *middleware) HttpMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer = &responseWriter{ctx.Writer}
		ctx.Next()
	}
}

type responseWriter struct {
	gin.ResponseWriter
}

func (w *responseWriter) Write(b []byte) (int, error) {
	fmt.Printf("writing response: %+v\n", string(b))
	cnt, err := w.ResponseWriter.Write(b)
	return cnt, err
}

func (w *responseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}
