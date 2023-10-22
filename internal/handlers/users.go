package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ShpullRequest/backend/internal/errs"
	"github.com/ShpullRequest/backend/internal/models"
	"github.com/ShpullRequest/backend/pkg/ip"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
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
			VkID:        int64(vkParams.VkUserID),
			SelectedGeo: "",
		})
		if err != nil {
			hs.logger.Error("Error create user", zap.Error(err))

			ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
			ctx.Abort()

			return
		}
		hs.logger.Debug("Success create user")
	}

	userGeo, err := ip.GetGeoByIP(ctx.ClientIP())
	if err != nil {
		hs.logger.Error("Error get user geo", zap.Error(err))
	}

	hs.logger.Debug("Success get me", zap.Any("User", user))
	ctx.JSON(http.StatusOK, models.NewResponse(struct {
		*models.User
		CurrentGeo string `json:"current_geo"`
	}{
		User:       user,
		CurrentGeo: userGeo,
	}))
	ctx.Abort()
}

func (hs *handlerService) GetUserByVkID(ctx *gin.Context) {
	var params struct {
		VkID string `uri:"vkId" binding:"required,numeric"`
	}

	if response, statusCode, err := hs.validateAndShouldBindURI(ctx, &params); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}

	vkID, _ := strconv.Atoi(params.VkID)

	user, err := hs.pg.GetUserByVkID(ctx, int64(vkID))
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
		PassedOnboarding bool    `json:"passed_onboarding,omitempty"`
		SelectedGeoLat   float64 `json:"selected_geo_lat" binding:"latitude"`
		SelectedGeoLot   float64 `json:"selected_geo_lot" binding:"longitude"`
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

	if !user.PassedOnboarding && params.PassedOnboarding {
		user.PassedOnboarding = true
	}
	if params.SelectedGeoLat != 0 && params.SelectedGeoLot != 0 {
		user.SelectedGeo = fmt.Sprintf("%f, %f", params.SelectedGeoLat, params.SelectedGeoLot)
	}

	if err = hs.pg.SaveUser(ctx, user); err != nil {
		hs.logger.Error("Error save user", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(user))
	ctx.Abort()
}
