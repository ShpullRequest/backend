package handlers

import (
	"database/sql"
	"errors"
	"github.com/ShpullRequest/backend/internal/errs"
	"github.com/ShpullRequest/backend/internal/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func (hs *handlerService) GetMe(ctx *gin.Context) {
	vkParams := hs.GetVKParams(ctx)

	var user *models.User
	user, err := hs.pg.GetUserByVkID(ctx, int64(vkParams.VkUserID))

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			hs.logger.Debug("Error get user", zap.Error(err))

			ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
			ctx.Abort()

			return
		}

		user, err = hs.pg.NewUser(ctx, models.User{
			VkID: int64(vkParams.VkUserID),
		})
		if err != nil {
			hs.logger.Error("Error create user", zap.Error(err))

			ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
			ctx.Abort()

			return
		}
		hs.logger.Debug("Success create user")
	}

	hs.logger.Debug("Success get me", zap.Any("User", user))
	ctx.JSON(http.StatusOK, models.NewResponse(user))
	ctx.Abort()
}

func (hs *handlerService) GetUserByVkID(ctx *gin.Context) {
	var params struct {
		VkID int64 `uri:"vkId" binding:"numeric"`
	}

	if response, statusCode, err := hs.validateAndShouldBindURI(ctx, &params); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}

	user, err := hs.pg.GetUserByVkID(ctx, params.VkID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	if !user.IsNil() {
		ctx.JSON(http.StatusOK, models.NewResponse(user))
	} else {
		ctx.JSON(http.StatusOK, models.NewResponse(nil))
	}

	ctx.Abort()
}

func (hs *handlerService) EditUser(ctx *gin.Context) {
	var params struct {
		PassedAppOnboarding    bool `json:"passed_app_onboarding,omitempty" binding:"required_without=PassedPrismaOnboarding"`
		PassedPrismaOnboarding bool `json:"passed_prisma_onboarding,omitempty" binding:"required_without=PassedAppOnboarding"`
	}

	if response, statusCode, err := hs.validateAndShouldBindJSON(ctx, &params); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}
	vkParams := hs.GetVKParams(ctx)

	user, err := hs.pg.GetUserByVkID(ctx, int64(vkParams.VkUserID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(errs.NewBadRequest("User not registered")))
		} else {
			ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		}

		ctx.Abort()
		return
	}

	if !user.PassedAppOnboarding && params.PassedAppOnboarding {
		user.PassedAppOnboarding = true
	}

	if !user.PassedPrismaOnboarding && params.PassedPrismaOnboarding {
		user.PassedPrismaOnboarding = true
	}

	if err = hs.pg.EditUser(ctx, user); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(user))
	ctx.Abort()
}
