package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/ShpullRequest/backend/internal/config"
	"github.com/ShpullRequest/backend/internal/errs"
	"github.com/ShpullRequest/backend/internal/models"
	"github.com/ShpullRequest/backend/pkg/vk/maps"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (hs *handlerService) GetReviewsPlace(ctx *gin.Context) {
	var params struct {
		PlaceID string `uri:"placeId" binding:"required,uuid"`
	}

	placeID, _ := uuid.Parse(params.PlaceID)
	reviews, err := hs.pg.GetReviewsPlace(ctx, placeID)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		hs.logger.Debug("Error get reviews", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	hs.logger.Debug("Success get reviews", zap.Any("Reviews", reviews))
	ctx.JSON(http.StatusOK, models.NewResponse(reviews))
	ctx.Abort()
}

func (hs *handlerService) NewReviewPlace(ctx *gin.Context) {
	vkParams := hs.GetVKParams(ctx)

	var params struct {
		PlaceID    string  `uri:"placeID" binding:"required,uuid"`
		ReviewText string  `uri:"review_text" binding:"required"`
		Stars      float64 `uri:"stars" binding:"required"`
	}

	if response, statusCode, err := hs.validateAndShouldBindJSON(ctx, &params); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}

	user, err := hs.pg.GetUserByVkID(ctx, int64(vkParams.VkUserID))
	if err != nil {
		hs.logger.Error("Error get user", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	placeID, _ := uuid.Parse(params.PlaceID)
	reviewPlace, err := hs.pg.NewReviewPlace(ctx, models.ReviewPlace{
		OwnerID:    user.ID,
		PlaceID:    placeID,
		ReviewText: params.ReviewText,
		Stars:      params.Stars,
		CreatedAt:  time.Now(),
	})

	if err != nil {
		hs.logger.Error("Error new review place", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(reviewPlace))
	ctx.Abort()
}

func (hs *handlerService) NewPlace(ctx *gin.Context) {
	var params struct {
		Name        string   `uri:"name" binding:"required"`
		Description string   `uri:"description" binding:"required"`
		Carousel    []string `uri:"carousel" binding:"required"`
		AddressLng  float64  `uri:"address_lng" binding:"required"`
		AddressLat  float64  `uri:"address_lat" binding:"required"`
	}

	if response, statusCode, err := hs.validateAndShouldBindJSON(ctx, &params); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}

	addressText, err := maps.New(config.Config).GetAddressByGeo(params.AddressLng, params.AddressLat)
	if err != nil {
		hs.logger.Error("Error get company by id", zap.Error(err))

		ctx.JSON(http.StatusBadGateway, models.NewErrorResponse(errs.NewBadGateway("Internal server error on vk maps")))
		ctx.Abort()

		return
	}

	place, err := hs.pg.NewPlace(ctx, models.Place{
		Name:        params.Name,
		Description: params.Description,
		Carousel:    params.Carousel,
		AddressText: addressText,
		AddressLng:  params.AddressLng,
		AddressLat:  params.AddressLat,
	})
	if err != nil {
		hs.logger.Error("Error new place", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(place))
	ctx.Abort()
}

func (hs *handlerService) EditPlace(ctx *gin.Context) {
	var params struct {
		PlaceID string `uri:"placeID" binding:"required,uuid"`
	}
	placeID, _ := uuid.Parse(params.PlaceID)

	if response, statusCode, err := hs.validateAndShouldBindJSON(ctx, &params); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}

	place, err := hs.pg.GetPlace(ctx, placeID)
	if err != nil {
		hs.logger.Error("Error get place", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	if err = hs.pg.SavePlace(ctx, place); err != nil {
		hs.logger.Error("Error save place", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}
}

func (hs *handlerService) GetPlace(ctx *gin.Context) {
	var params struct {
		PlaceID string `uri:"placeID" binding:"required,uuid"`
	}

	if response, statusCode, err := hs.validateAndShouldBindURI(ctx, &params); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}

	placeID, _ := uuid.Parse(params.PlaceID)

	place, err := hs.pg.GetPlace(ctx, placeID)
	if err != nil {
		hs.logger.Error("Error get place", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(place))
	ctx.Abort()
}

func (hs *handlerService) GetAllPlaces(ctx *gin.Context) {
	places, err := hs.pg.GetAllPlaces(ctx)
	if err != nil {
		hs.logger.Error("Error get all places", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(places))
	ctx.Abort()
}
