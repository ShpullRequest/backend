package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (ms *middlewareService) Logger(ctx *gin.Context) {
	startRequest := time.Now()

	ctx.Next()

	timeSinceRequest := time.Since(startRequest)
	method := ctx.Request.Method
	path := ctx.Request.URL.Path
	statusCode := ctx.Writer.Status()
	size := ctx.Writer.Size()

	ms.logger.Info(
		"Request",
		zap.Int64("Duration", timeSinceRequest.Milliseconds()),
		zap.String("Method", method),
		zap.String("Path", path),
		zap.Int("StatusCode", statusCode),
		zap.Int("Size", size),
	)
}
