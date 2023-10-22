package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (ms *middlewareService) Cors(ctx *gin.Context) {
	if ctx.Request.Method == http.MethodOptions {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Headers", "*")
		ctx.Header("Access-Control-Allow-Methods", "*")

		ctx.Status(http.StatusOK)
		ctx.Abort()
	}
}
