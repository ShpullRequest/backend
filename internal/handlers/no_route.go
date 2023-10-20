package handlers

import (
	"net/http"

	"github.com/ShpullRequest/backend/internal/errs"
	"github.com/gin-gonic/gin"
)

func (hs *handlerService) NoRoute(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, errs.NewNotFound("Invalid method path"))
	ctx.Abort()
}
