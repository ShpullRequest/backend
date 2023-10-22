package handlers

import (
	"database/sql"
	"errors"
	"github.com/jackc/pgerrcode"
	"math"
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

// NewPlace
// @Summary Добавить новое место
// @Description Создает новое место.
// @ID create-place
// @Accept json
// @Produce json
// @Param Authorization header string true "Строка авторизации"
// @Param name body string true "Название места"
// @Param description body string true "Описание места"
// @Param carousel body []string true "Список изображений для карусели"
// @Param address_lng body float64 true "Долгота местоположения"
// @Param address_lat body float64 true "Широта местоположения"
// @Success 200 {object} models.Place
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /places [post]
func (hs *handlerService) NewPlace(ctx *gin.Context) {
	var params struct {
		Name        string   `json:"name" binding:"required,min=6"`
		Description string   `json:"description" binding:"required,min=10"`
		Carousel    []string `json:"carousel" binding:"required"`
		AddressLng  float64  `json:"address_lng" binding:"required"`
		AddressLat  float64  `json:"address_lat" binding:"required"`
	}

	if response, statusCode, err := hs.validateAndShouldBindJSON(ctx, &params); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}

	addressText, err := maps.New(config.Config).GetAddressByGeo(params.AddressLng, params.AddressLat)
	if err != nil {
		hs.logger.Error("Error get address", zap.Error(err))

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

// EditPlace
// @Summary Редактировать место
// @Description Редактирует существующее место.
// @ID edit-place
// @Accept json
// @Produce json
// @Param Authorization header string true "Строка авторизации"
// @Param placeId path string true "Уникальный идентификатор места (в формате UUID)"
// @Success 200 {object} models.Place
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /places/{placeID} [patch]
func (hs *handlerService) EditPlace(ctx *gin.Context) {
	var paramsURI struct {
		PlaceID string `uri:"placeId" binding:"required,uuid"`
	}

	if response, statusCode, err := hs.validateAndShouldBindURI(ctx, &paramsURI); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}

	var params struct {
		Name        string   `json:"name" binding:"omitempty,min=6"`
		Description string   `json:"description" binding:"omitempty,min=10"`
		Carousel    []string `json:"carousel"`
		AddressLng  float64  `json:"address_lng"`
		AddressLat  float64  `json:"address_lat"`
	}

	if response, statusCode, err := hs.validateAndShouldBindJSON(ctx, &params); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}

	placeID, _ := uuid.Parse(paramsURI.PlaceID)
	place, err := hs.pg.GetPlace(ctx, placeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, models.NewErrorResponse(errs.NewNotFound("Place not found")))
		} else {
			hs.logger.Error("Error get place", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		}
		ctx.Abort()

		return
	}

	if params.Name != "" {
		place.Name = params.Name
	}
	if params.Description != "" {
		place.Description = params.Description
	}
	if len(params.Carousel) > 0 {
		place.Carousel = params.Carousel
	}
	if params.AddressLng != 0 && params.AddressLat != 0 {
		address, err := maps.New(config.Config).GetAddressByGeo(params.AddressLng, params.AddressLat)
		if err != nil {
			hs.logger.Error("Error get address", zap.Error(err))

			ctx.JSON(http.StatusBadGateway, models.NewErrorResponse(errs.NewBadGateway("Internal server error on vk maps")))
			ctx.Abort()

			return
		}

		place.AddressText = address
		place.AddressLng = params.AddressLng
		place.AddressLat = params.AddressLat
	}

	if err = hs.pg.SavePlace(ctx, place); err != nil {
		hs.logger.Error("Error save place", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(place))
}

// GetPlace
// @Summary Получить место
// @Description Возвращает информацию о указанном месте.
// @ID get-place
// @Accept json
// @Produce json
// @Param Authorization header string true "Строка авторизации"
// @Param placeId path string true "Уникальный идентификатор места (в формате UUID)"
// @Success 200 {object} models.Place
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /places/{placeID} [get]
func (hs *handlerService) GetPlace(ctx *gin.Context) {
	var params struct {
		PlaceID string `uri:"placeId" binding:"required,uuid"`
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

// SearchPlaces
// @Summary Поиск мест
// @Description Ищет места по заданному запросу.
// @ID search-routes
// @Accept json
// @Produce json
// @Param Authorization header string true "Строка авторизации"
// @Param query path string true "Запрос для поиска мест (минимум 2 символа)"
// @Success 200 {object} []models.Place
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /routes/search/{query} [get]
func (hs *handlerService) SearchPlaces(ctx *gin.Context) {
	var params struct {
		Query string `uri:"query" binding:"required,min=2"`
	}

	if response, statusCode, err := hs.validateAndShouldBindURI(ctx, &params); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}

	places, err := hs.pg.SearchPlace(ctx, params.Query)
	if err != nil {
		hs.logger.Error("Error search routes", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(places))
	ctx.Abort()
}

// GetAllPlaces
// @Summary Получить все места
// @Description Возвращает список всех мест.
// @ID get-all-places
// @Accept json
// @Produce json
// @Param Authorization header string true "Строка авторизации"
// @Success 200 {object} []models.Place
// @Failure 500 {object} models.ErrorResponse
// @Router /places [get]
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

// NewReviewPlace
// @Summary Добавить новый отзыв о месте
// @Description Создает новый отзыв о указанном месте.
// @ID create-place-review
// @Accept json
// @Produce json
// @Param Authorization header string true "Строка авторизации"
// @Param placeId path string true "Уникальный идентификатор места (в формате UUID)"
// @Param review_text body string true "Текст отзыва (минимум 6 символов)"
// @Param stars body float64 true "Оценка места (от 1 до 5)"
// @Success 200 {object} models.ReviewPlace
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /places/{placeID}/reviews [post]
func (hs *handlerService) NewReviewPlace(ctx *gin.Context) {
	var paramsURI struct {
		PlaceID string `uri:"placeId" binding:"required,uuid"`
	}

	if response, statusCode, err := hs.validateAndShouldBindURI(ctx, &paramsURI); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}

	var params struct {
		ReviewText string  `json:"review_text" binding:"required,min=6"`
		Stars      float64 `json:"stars" binding:"required"`
	}

	if response, statusCode, err := hs.validateAndShouldBindJSON(ctx, &params); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}
	vkParams := hs.GetVKParams(ctx)

	user, err := hs.pg.GetUserByVkID(ctx, int64(vkParams.VkUserID))
	if err != nil {
		hs.logger.Error("Error get user", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	if params.Stars < 1 || params.Stars > 5 {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(errs.NewBadRequest("Field validation for \"Stars\" failed on the 'min=1,max=5' tag.")))
		ctx.Abort()

		return
	} else {
		params.Stars = math.Round(params.Stars)
	}

	placeID, _ := uuid.Parse(paramsURI.PlaceID)
	reviewPlace, err := hs.pg.NewReviewPlace(ctx, models.ReviewPlace{
		OwnerID:    user.ID,
		PlaceID:    placeID,
		ReviewText: params.ReviewText,
		Stars:      params.Stars,
		CreatedAt:  time.Now(),
	})

	if err != nil {
		if hs.pg.IsError(pgerrcode.IsIntegrityConstraintViolation, err) {
			ctx.JSON(http.StatusConflict, models.NewErrorResponse(errs.NewConflict("You have already added a review to this place")))
		} else {
			hs.logger.Error("Error new review event", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		}
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(reviewPlace))
	ctx.Abort()
}

// EditReviewPlace
// @Summary Редактировать отзыв о месте
// @Description Редактирует существующий отзыв о месте с указанными параметрами.
// @ID edit-review-route
// @Accept json
// @Produce json
// @Param Authorization header string true "Строка авторизации"
// @Param routeId path string true "Уникальный идентификатор места"
// @Param review_text body string false "Текст отзыва (минимум 6 символов, опционально)"
// @Param stars body number false "Оценка (от 1 до 5, опционально)"
// @Success 200 {object} models.ReviewPlace
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /routes/{routeId}/reviews [patch]
func (hs *handlerService) EditReviewPlace(ctx *gin.Context) {
	var paramsURI struct {
		RouteID string `uri:"routeId" binding:"required,uuid"`
	}

	if response, statusCode, err := hs.validateAndShouldBindURI(ctx, &paramsURI); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}

	var params struct {
		ReviewText string  `json:"review_text"`
		Stars      float64 `json:"stars"`
	}

	if response, statusCode, err := hs.validateAndShouldBindJSON(ctx, &params); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}
	vkParams := hs.GetVKParams(ctx)

	user, err := hs.pg.GetUserByVkID(ctx, int64(vkParams.VkUserID))
	if err != nil {
		hs.logger.Debug("Error get user", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	routeID, _ := uuid.Parse(paramsURI.RouteID)
	reviewRoute, err := hs.pg.GetReviewRoute(ctx, user.ID, routeID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, models.NewErrorResponse(errs.NewNotFound("Review route not found")))
		} else {
			hs.logger.Error("Error get event", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		}
		ctx.Abort()

		return
	}

	if params.ReviewText != "" {
		if len(params.ReviewText) < 6 {
			ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(errs.NewBadRequest("Field validation for \"ReviewText\" failed on the 'min=6' tag.")))
			ctx.Abort()

			return
		}

		reviewRoute.ReviewText = params.ReviewText
	}
	if params.Stars != 0 {
		if params.Stars < 1 || params.Stars > 5 {
			ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(errs.NewBadRequest("Field validation for \"Stars\" failed on the 'min=1,max=5' tag.")))
			ctx.Abort()

			return
		} else {
			reviewRoute.Stars = math.Round(params.Stars)
		}
	}

	if err = hs.pg.SaveReviewRoute(ctx, reviewRoute); err != nil {
		hs.logger.Debug("Error save review route", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(reviewRoute))
	ctx.Abort()
}

// GetReviewsPlace
// @Summary Получить отзывы о месте
// @Description Возвращает список всех отзывов о указанном месте.
// @ID get-place-reviews
// @Accept json
// @Produce json
// @Param Authorization header string true "Строка авторизации"
// @Param placeId path string true "Уникальный идентификатор места (в формате UUID)"
// @Success 200 {object} []models.ReviewPlace
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /places/{placeId}/reviews [get]
func (hs *handlerService) GetReviewsPlace(ctx *gin.Context) {
	var params struct {
		PlaceID string `uri:"placeId" binding:"required,uuid"`
	}

	if response, statusCode, err := hs.validateAndShouldBindURI(ctx, &params); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}

	placeID, _ := uuid.Parse(params.PlaceID)
	places, err := hs.pg.GetReviewsEvent(ctx, placeID)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		hs.logger.Debug("Error get reviews", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(places))
	ctx.Abort()
}
