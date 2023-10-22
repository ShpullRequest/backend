package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/ShpullRequest/backend/internal/errs"
	"github.com/ShpullRequest/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (hs *handlerService) NewAchievement(ctx *gin.Context) {
	var params struct {
		Name        string `json:"name" binding:"required,min=6"`
		Description string `json:"description" binding:"required,min=10"`
		Icon        string `json:"icon" binding:"required,url"`
		Coins       int    `json:"coins" binding:"required"`
	}

	if response, statusCode, err := hs.validateAndShouldBindJSON(ctx, &params); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}

	user, err := hs.pg.GetUserByVkID(ctx, int64(hs.GetVKParams(ctx).VkUserID))

	if err != nil {
		hs.logger.Error("Error get user by vk id", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	if !user.IsAdmin {
		ctx.JSON(http.StatusForbidden, models.NewErrorResponse(errs.NewForbidden("You don't have access to this method")))
		ctx.Abort()

		return
	}

	achievement, err := hs.pg.NewAchievement(ctx, models.Achievements{
		Name:        params.Name,
		Description: params.Description,
		Icon:        params.Icon,
		Coins:       params.Coins,
	})

	if err != nil {
		hs.logger.Error("Error new achievement", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(achievement))
	ctx.Abort()
}

func (hs *handlerService) GetAchievement(ctx *gin.Context) {
	var params struct {
		AchievementID string `uri:"achievementId" blinding:"required,uuid"`
	}

	if response, statusCode, err := hs.validateAndShouldBindURI(ctx, &params); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}

	achievementID, _ := uuid.Parse(params.AchievementID)
	achievement, err := hs.pg.GetAchievementByID(ctx, achievementID)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		hs.logger.Error("Error get achievement", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(achievement))
	ctx.Abort()
}

func (hs *handlerService) EditAchievement(ctx *gin.Context) {
	var params struct {
		AchievementID string `uri:"achievementId" binding:"required,uuid"`
	}

	if response, statusCode, err := hs.validateAndShouldBindURI(ctx, &params); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}

	user, err := hs.pg.GetUserByVkID(ctx, int64(hs.GetVKParams(ctx).VkUserID))
	if err != nil {
		hs.logger.Error("Error get user by vk id", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	if !user.IsAdmin {
		ctx.JSON(http.StatusForbidden, models.NewErrorResponse(errs.NewForbidden("You don't have access to this method")))
		ctx.Abort()

		return
	}

	achievementID, _ := uuid.Parse(params.AchievementID)
	achievement, err := hs.pg.GetAchievementByID(ctx, achievementID)
	if err != nil {
		hs.logger.Error("Error get achievement", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	if err = hs.pg.SaveAchievement(ctx, achievement); err != nil {
		hs.logger.Error("Error save achievement", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}
}

func (hs *handlerService) GetAllAchievements(ctx *gin.Context) {
	achievements, err := hs.pg.GetAllAchievements(ctx)
	if err != nil {
		hs.logger.Error("Error get all achievements", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(achievements))
	ctx.Abort()
}
