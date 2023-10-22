package handlers

import (
	"database/sql"
	"errors"
	"github.com/ShpullRequest/backend/internal/config"
	"github.com/ShpullRequest/backend/internal/errs"
	"github.com/ShpullRequest/backend/internal/models"
	"github.com/ShpullRequest/backend/pkg/vk/maps"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"go.uber.org/zap"
	"math"
	"net/http"
	"time"
)

// NewEvent
// @Summary Создать новое событие
// @Description Создает новое событие с указанными параметрами.
// @ID create-event
// @Accept json
// @Produce json
// @Param company_id body string false "Уникальный идентификатор компании (в формате UUID)"
// @Param name body string true "Название события (минимум 6 символов)"
// @Param description body string true "Описание события (минимум 10 символов)"
// @Param carousel body []string true "Массив ссылок на изображения для карусели события"
// @Param tags body []string true "Массив тегов для события"
// @Param icon body string true "Ссылка на иконку события (должна быть валидной URL)"
// @Param start_time body string true "Дата и время начала события (в формате 2006-01-02T15:04:05Z07:00)"
// @Param address_lng body float64 true "Долгота местоположения события"
// @Param address_lat body float64 true "Широта местоположения события"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /events [post]
func (hs *handlerService) NewEvent(ctx *gin.Context) {
	var params struct {
		// Нельзя парсить сразу в uuid.UUID, потому что это не поддерживает нормально gin
		CompanyID   string   `json:"company_id" binding:"omitempty,uuid"`
		Name        string   `json:"name" binding:"required,min=6"`
		Description string   `json:"description" binding:"required,min=10"`
		Carousel    []string `json:"carousel" binding:"required"`
		Tags        []string `json:"tags" binding:"required"`
		Icon        string   `json:"icon" binding:"required,url"`
		StartTime   string   `json:"start_time" binding:"required"`
		AddressLng  float64  `json:"address_lng" binding:"required,longitude"`
		AddressLat  float64  `json:"address_lat" binding:"required,latitude"`
	}

	if response, statusCode, err := hs.validateAndShouldBindJSON(ctx, &params); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}

	startTime, err := time.Parse("2006-01-02T15:04:05Z07:00", params.StartTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(errs.NewBadRequest("Field validation for \"StartTime\" failed on the 'timezone=2006-01-02T15:04:05Z07:00' tag.")))
		ctx.Abort()

		return
	}

	vkParams := hs.GetVKParams(ctx)

	user, err := hs.pg.GetUserByVkID(ctx, int64(vkParams.VkUserID))
	if err != nil {
		hs.logger.Error("Error get user by vk id", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	var companyID *uuid.UUID = nil
	paramCompanyID, err := uuid.Parse(params.CompanyID)

	if err != nil {
		if !user.IsAdmin {
			ctx.JSON(http.StatusForbidden, models.NewErrorResponse(errs.NewForbidden("You don't have access to this method")))
			ctx.Abort()

			return
		}
	} else {
		company, err := hs.pg.GetCompanyByID(ctx, paramCompanyID)
		if err != nil {
			hs.logger.Error("Error get company by id", zap.Error(err))

			ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
			ctx.Abort()

			return
		} else if company.UserID != user.ID {
			ctx.JSON(http.StatusForbidden, models.NewErrorResponse(errs.NewForbidden("You can't create an event on behalf of this company")))
			ctx.Abort()

			return
		}

		companyID = &company.ID
	}

	address, err := maps.New(config.Config).GetAddressByGeo(params.AddressLng, params.AddressLat)
	if err != nil {
		hs.logger.Error("Error get company by id", zap.Error(err))

		ctx.JSON(http.StatusBadGateway, models.NewErrorResponse(errs.NewBadGateway("Internal server error on vk maps")))
		ctx.Abort()

		return
	}

	event, err := hs.pg.NewEvent(ctx, models.Event{
		CompanyID:   companyID,
		Name:        params.Name,
		Description: params.Description,
		Carousel:    params.Carousel,
		Tags:        params.Tags,
		Icon:        params.Icon,
		StartTime:   startTime,
		AddressText: address,
		AddressLng:  params.AddressLng,
		AddressLat:  params.AddressLat,
	})
	if err != nil {
		hs.logger.Error("Error new event", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(event))
	ctx.Abort()
}

// GetEvent
// @Summary Получить информацию о событии
// @Description Возвращает информацию о конкретном событии по его ID.
// @ID get-event
// @Accept json
// @Produce json
// @Param eventId path string true "Уникальный идентификатор события (в формате UUID)"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /events/{eventId} [get]
func (hs *handlerService) GetEvent(ctx *gin.Context) {
	var params struct {
		EventID string `uri:"eventId" binding:"required,uuid"`
	}

	if response, statusCode, err := hs.validateAndShouldBindURI(ctx, &params); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}

	eventID, _ := uuid.Parse(params.EventID)
	event, err := hs.pg.GetEvent(ctx, eventID)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		hs.logger.Error("Error get event", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	if !event.IsNil() {
		ctx.JSON(http.StatusOK, models.NewResponse(event))
	} else {
		ctx.JSON(http.StatusOK, models.NewResponse(nil))
	}
	ctx.Abort()
}

// EditEvent
// @Summary Редактировать событие
// @Description Редактирует информацию о существующем событии.
// @ID edit-event
// @Accept json
// @Produce json
// @Param eventId path string true "Уникальный идентификатор события (в формате UUID)"
// @Param name body string false "Новое название события (минимум 6 символов)"
// @Param description body string false "Новое описание события (минимум 10 символов)"
// @Param carousel body []string false "Новый массив ссылок на изображения для карусели события"
// @Param tags body []string false "Новый массив тегов для события"
// @Param icon body string false "Новая ссылка на иконку события (должна быть валидной URL)"
// @Param start_time body string false "Новая дата и время начала события (в формате 2006-01-02T15:04:05Z07:00)"
// @Param address_lng body float64 false "Новая долгота местоположения события"
// @Param address_lat body float64 false "Новая широта местоположения события"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /events/{eventId} [patch]
func (hs *handlerService) EditEvent(ctx *gin.Context) {
	var paramsURI struct {
		EventID string `uri:"eventId" binding:"required,uuid"`
	}

	if response, statusCode, err := hs.validateAndShouldBindURI(ctx, &paramsURI); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}

	var params struct {
		Name        string   `json:"name" binding:"omitempty,min=6"`
		Description string   `json:"description" binding:"omitempty,min=10"`
		Carousel    []string `json:"carousel" binding:"omitempty"`
		Tags        []string `json:"tags" binding:"required"`
		Icon        string   `json:"icon" binding:"omitempty,url"`
		StartTime   string   `json:"start_time" binding:"omitempty"`
		AddressLng  float64  `json:"address_lng" binding:"omitempty,longitude"`
		AddressLat  float64  `json:"address_lat" binding:"omitempty,latitude"`
	}

	if response, statusCode, err := hs.validateAndShouldBindJSON(ctx, &params); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}
	vkParams := hs.GetVKParams(ctx)

	eventID, _ := uuid.Parse(paramsURI.EventID)
	event, err := hs.pg.GetEvent(ctx, eventID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, models.NewErrorResponse(errs.NewNotFound("Event not found")))
		} else {
			hs.logger.Error("Error get event", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		}
		ctx.Abort()

		return
	}

	user, err := hs.pg.GetUserByVkID(ctx, int64(vkParams.VkUserID))
	if err != nil {
		hs.logger.Error("Error get user by vk id", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	if event.CompanyID == nil {
		if !user.IsAdmin {
			ctx.JSON(http.StatusForbidden, models.NewErrorResponse(errs.NewForbidden("You don't have access to this method")))
			ctx.Abort()

			return
		}
	} else {
		company, err := hs.pg.GetCompanyByID(ctx, *event.CompanyID)
		if err != nil {
			hs.logger.Error("Error get company by id", zap.Error(err))

			ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
			ctx.Abort()

			return
		}

		if company.UserID != user.ID {
			ctx.JSON(http.StatusForbidden, models.NewErrorResponse(errs.NewForbidden("You can't create an event on behalf of this company")))
			ctx.Abort()

			return
		}
	}

	if params.Name != "" {
		event.Name = params.Name
	}
	if params.Description != "" {
		event.Description = params.Description
	}
	if len(params.Carousel) > 0 {
		event.Carousel = params.Carousel
	}
	if len(params.Tags) > 0 {
		event.Tags = params.Tags
	}
	if params.Icon != "" {
		event.Icon = params.Icon
	}
	if params.StartTime != "" {
		startTime, err := time.Parse("2006-01-02T15:04:05Z07:00", params.StartTime)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(errs.NewBadRequest("Field validation for \"StartTime\" failed on the 'timezone=2006-01-02T15:04:05Z07:00' tag.")))
			ctx.Abort()

			return
		}
		event.StartTime = startTime
	}
	if params.AddressLng != 0 && params.AddressLat != 0 {
		address, err := maps.New(config.Config).GetAddressByGeo(params.AddressLng, params.AddressLat)
		if err != nil {
			hs.logger.Error("Error get company by id", zap.Error(err))

			ctx.JSON(http.StatusBadGateway, models.NewErrorResponse(errs.NewBadGateway("Internal server error on vk maps")))
			ctx.Abort()

			return
		}

		event.AddressText = address
		event.AddressLng = params.AddressLng
		event.AddressLat = params.AddressLat
	}

	if err = hs.pg.SaveEvent(ctx, event); err != nil {
		hs.logger.Error("Error save event", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(event))
	ctx.Abort()
}

// GetCompanyEvents
// @Summary Получить все события компании
// @Description Возвращает список всех событий, принадлежащих конкретной компании.
// @ID get-company-events
// @Accept json
// @Produce json
// @Param companyId path string true "Уникальный идентификатор компании (в формате UUID)"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /companies/{companyId}/events [get]
func (hs *handlerService) GetCompanyEvents(ctx *gin.Context) {
	var params struct {
		CompanyID string `uri:"companyId" binding:"required,uuid"`
	}

	if response, statusCode, err := hs.validateAndShouldBindURI(ctx, &params); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}

	companyID, _ := uuid.Parse(params.CompanyID)
	events, err := hs.pg.GetAllEventsByCompanyID(ctx, companyID)

	if err != nil {
		hs.logger.Error("Error get all events by company id", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(events))
	ctx.Abort()
}

// SearchEvents
// @Summary Поиск событий
// @Description Ищет события по заданному запросу.
// @ID search-events
// @Accept json
// @Produce json
// @Param query path string true "Поисковый запрос (минимум 2 символа)"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /events/search/{query} [get]
func (hs *handlerService) SearchEvents(ctx *gin.Context) {
	var params struct {
		Query string `uri:"query" binding:"required,min=2"`
	}

	if response, statusCode, err := hs.validateAndShouldBindURI(ctx, &params); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}

	events, err := hs.pg.SearchEvents(ctx, params.Query)
	if err != nil {
		hs.logger.Error("Error search events", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(events))
	ctx.Abort()
}

// GetAllEvents
// @Summary Получить все события
// @Description Возвращает список всех событий в системе.
// @ID get-all-events
// @Accept json
// @Produce json
// @Success 200 {object} models.Response
// @Failure 500 {object} models.ErrorResponse
// @Router /events [get]
func (hs *handlerService) GetAllEvents(ctx *gin.Context) {
	events, err := hs.pg.GetAllEvents(ctx)
	if err != nil {
		hs.logger.Error("Error get all events", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(events))
	ctx.Abort()
}

// NewReviewEvent
// @Summary Добавить новый отзыв к событию
// @Description Создает новый отзыв к указанному событию.
// @ID create-event-review
// @Accept json
// @Produce json
// @Param eventId path string true "Уникальный идентификатор события (в формате UUID)"
// @Param review_text body string true "Текст отзыва (минимум 6 символов)"
// @Param stars body float64 true "Оценка события (от 1 до 5)"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /events/{eventId}/reviews [post]
func (hs *handlerService) NewReviewEvent(ctx *gin.Context) {
	var paramsURI struct {
		EventID string `uri:"eventId" binding:"required,uuid"`
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

	eventID, _ := uuid.Parse(paramsURI.EventID)
	reviewEvent, err := hs.pg.NewReviewEvent(ctx, models.ReviewEvent{
		OwnerID:    user.ID,
		EventID:    eventID,
		ReviewText: params.ReviewText,
		Stars:      params.Stars,
		CreatedAt:  time.Now(),
	})

	if err != nil {
		if hs.pg.IsError(pgerrcode.IsIntegrityConstraintViolation, err) {
			ctx.JSON(http.StatusConflict, models.NewErrorResponse(errs.NewConflict("You have already added a review to this event")))
		} else {
			hs.logger.Error("Error new review event", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		}
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(reviewEvent))
	ctx.Abort()
}

// EditReviewsEvent
// @Summary Редактировать отзыв к событию
// @Description Редактирует существующий отзыв к событию.
// @ID edit-event-review
// @Accept json
// @Produce json
// @Param eventId path string true "Уникальный идентификатор события (в формате UUID)"
// @Param review_text body string false "Новый текст отзыва (минимум 6 символов)"
// @Param stars body float64 false "Новая оценка события (от 1 до 5)"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /events/{eventId}/reviews [patch]
func (hs *handlerService) EditReviewsEvent(ctx *gin.Context) {
	var paramsURI struct {
		EventID string `uri:"eventId" binding:"required,uuid"`
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

	eventID, _ := uuid.Parse(paramsURI.EventID)
	reviewEvent, err := hs.pg.GetReviewEvent(ctx, user.ID, eventID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, models.NewErrorResponse(errs.NewNotFound("Review event not found")))
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

		reviewEvent.ReviewText = params.ReviewText
	}
	if params.Stars != 0 {
		if params.Stars < 1 || params.Stars > 5 {
			ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(errs.NewBadRequest("Field validation for \"Stars\" failed on the 'min=1,max=5' tag.")))
			ctx.Abort()

			return
		} else {
			reviewEvent.Stars = math.Round(params.Stars)
		}
	}

	if err = hs.pg.SaveReviewEvent(ctx, reviewEvent); err != nil {
		hs.logger.Debug("Error save review event", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(reviewEvent))
	ctx.Abort()
}

// GetReviewsEvent
// @Summary Получить отзывы к событию
// @Description Возвращает список всех отзывов к указанному событию.
// @ID get-event-reviews
// @Accept json
// @Produce json
// @Param eventId path string true "Уникальный идентификатор события (в формате UUID)"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /events/{eventId}/reviews [get]
func (hs *handlerService) GetReviewsEvent(ctx *gin.Context) {
	var params struct {
		EventID string `uri:"eventId" binding:"required,uuid"`
	}

	if response, statusCode, err := hs.validateAndShouldBindURI(ctx, &params); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}

	eventID, _ := uuid.Parse(params.EventID)
	reviews, err := hs.pg.GetReviewsEvent(ctx, eventID)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		hs.logger.Debug("Error get reviews", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(reviews))
	ctx.Abort()
}
