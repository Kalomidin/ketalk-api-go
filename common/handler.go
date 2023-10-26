package common

import (
	"ketalk-api/common/response"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerFunc func(ctx *gin.Context, req *http.Request) (interface{}, error)

func GenericHandler(h HandlerFunc) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		val := ctx.Request
		resp, err := h(ctx, val)
		var sendErr error
		if err != nil {
			log.Printf("failed with error: %v\n", err)
			sendErr = response.NewError(err, http.StatusInternalServerError).Send(ctx.Writer)
		} else {
			sendErr = response.NewSuccess(resp, http.StatusOK).Send(ctx.Writer)
		}
		if sendErr != nil {
			log.Printf("failed to send response: %v", sendErr)
		}
	}
}
