package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (ms *middlewareService) Cors(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "*")

	if ctx.Request.Method == http.MethodOptions {
		ctx.Status(http.StatusOK)
		ctx.Abort()
	}
}
