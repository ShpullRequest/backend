package middlewares

import "github.com/gin-gonic/gin"

func (ms *middlewareService) Cors(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "Authorization")
}
