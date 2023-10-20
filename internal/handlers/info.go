package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (hs *handlerService) Info(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
	ctx.Abort()
}
