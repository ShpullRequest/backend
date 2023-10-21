package handlers

import (
	"database/sql"
	"errors"
	"github.com/ShpullRequest/backend/internal/errs"
	"github.com/ShpullRequest/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func (hs *handlerService) NewEvent(ctx *gin.Context) {
	var params struct {
		CompanyID   uuid.UUID `json:"company_id,omitempty" binding:"omitempty,uuid"`
		Name        string    `json:"name" binding:"required,min=6"`
		Description string    `json:"description" binding:"required,min=10"`
		Carousel    []string  `json:"carousel" binding:"required"`
		Icon        string    `json:"icon" binding:"required,url"`
		StartTime   string    `json:"start_time" binding:"required"`
		AddressLng  float64   `json:"address_lng" binding:"required,longitude"`
		AddressLat  float64   `json:"address_lat" binding:"required,latitude"`
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

	event, err := hs.pg.NewEvent(ctx, models.Event{
		CompanyID:   &params.CompanyID,
		Name:        params.Name,
		Description: params.Description,
		Carousel:    params.Carousel,
		Icon:        params.Icon,
		StartTime:   startTime,
		AddressText: "", // TODO: by vk maps api
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
