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
	"go.uber.org/zap"
	"net/http"
	"time"
)

func (hs *handlerService) NewEvent(ctx *gin.Context) {
	var params struct {
		// Нельзя парсить сразу в uuid.UUID, потому что это не поддерживает нормально gin
		CompanyID   string   `json:"company_id" binding:"omitempty,uuid"`
		Name        string   `json:"name" binding:"required,min=6"`
		Description string   `json:"description" binding:"required,min=10"`
		Carousel    []string `json:"carousel" binding:"required"`
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

func (hs *handlerService) SaveEvent(ctx *gin.Context) {
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
