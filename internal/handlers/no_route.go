package handlers

import (
	"github.com/ShpullRequest/backend/internal/models"
	"net/http"

	"github.com/ShpullRequest/backend/internal/errs"
	"github.com/gin-gonic/gin"
)

func (hs *handlerService) NoRoute(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, models.NewErrorResponse(errs.NewNotFound("Invalid method path")))
	ctx.Abort()
}
