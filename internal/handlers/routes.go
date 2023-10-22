package handlers

import (
	"database/sql"
	"errors"
	"github.com/ShpullRequest/backend/internal/errs"
	"github.com/ShpullRequest/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"go.uber.org/zap"
	"math"
	"net/http"
	"time"
)

func (hs *handlerService) NewRoute(ctx *gin.Context) {
	var params struct {
		CompanyID   string   `json:"company_id" binding:"omitempty,uuid"`
		Name        string   `json:"name" binding:"min=6"`
		Description string   `json:"description" binding:"min=10"`
		Places      []string `json:"places,omitempty"`
		Events      []string `json:"events,omitempty"`
	}

	if response, statusCode, err := hs.validateAndShouldBindJSON(ctx, &params); err != nil {
		ctx.JSON(statusCode, response)
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

	route, err := hs.pg.NewRoute(ctx, models.Route{
		CompanyID:   companyID,
		Name:        params.Name,
		Description: params.Description,
		Places:      params.Places,
		Events:      params.Events,
	})
	if err != nil {
		hs.logger.Error("Error new route", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(route))
	ctx.Abort()
}

func (hs *handlerService) EditRoute(ctx *gin.Context) {
	var paramsURI struct {
		RouteID string `uri:"routeId" binding:"required,uuid"`
	}

	if response, statusCode, err := hs.validateAndShouldBindURI(ctx, &paramsURI); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}

	var params struct {
		Name        string   `json:"name" binding:"omitempty,min=6"`
		Description string   `json:"description" binding:"omitempty,min=10"`
		Places      []string `json:"places,omitempty"`
		Events      []string `json:"events,omitempty"`
	}

	if response, statusCode, err := hs.validateAndShouldBindJSON(ctx, &params); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}
	vkParams := hs.GetVKParams(ctx)

	routeID, _ := uuid.Parse(paramsURI.RouteID)
	route, err := hs.pg.GetRoute(ctx, routeID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, models.NewErrorResponse(errs.NewNotFound("Route not found")))
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

	if route.CompanyID == nil {
		if !user.IsAdmin {
			ctx.JSON(http.StatusForbidden, models.NewErrorResponse(errs.NewForbidden("You don't have access to this method")))
			ctx.Abort()

			return
		}
	} else {
		company, err := hs.pg.GetCompanyByID(ctx, *route.CompanyID)
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
		route.Name = params.Name
	}
	if params.Description != "" {
		route.Description = params.Description
	}
	if len(params.Places) > 0 {
		route.Places = params.Places
	}
	if len(params.Events) > 0 {
		route.Events = params.Events
	}

	if err = hs.pg.SaveRoute(ctx, route); err != nil {
		hs.logger.Error("Error save route", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(route))
	ctx.Abort()
}

func (hs *handlerService) GetRoute(ctx *gin.Context) {
	var params struct {
		RouteID string `uri:"routeId" binding:"required,uuid"`
	}

	if response, statusCode, err := hs.validateAndShouldBindURI(ctx, &params); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}

	routeID, _ := uuid.Parse(params.RouteID)
	route, err := hs.pg.GetRoute(ctx, routeID)

	if err != nil {
		hs.logger.Error("Error get route", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(route))
	ctx.Abort()
}

func (hs *handlerService) SearchRoutes(ctx *gin.Context) {
	var params struct {
		Query string `uri:"query" binding:"required,min=2"`
	}

	if response, statusCode, err := hs.validateAndShouldBindURI(ctx, &params); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}

	routes, err := hs.pg.SearchRoutes(ctx, params.Query)
	if err != nil {
		hs.logger.Error("Error search routes", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(routes))
	ctx.Abort()
}

func (hs *handlerService) GetCompanyRoutes(ctx *gin.Context) {
	var params struct {
		CompanyID string `uri:"companyId" binding:"required,uuid"`
	}

	if response, statusCode, err := hs.validateAndShouldBindURI(ctx, &params); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}

	companyID, _ := uuid.Parse(params.CompanyID)
	routes, err := hs.pg.GetAllRoutesByCompanyID(ctx, companyID)

	if err != nil {
		hs.logger.Error("Error get all routes by company id", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(routes))
	ctx.Abort()
}

func (hs *handlerService) GetAllRoutes(ctx *gin.Context) {
	routes, err := hs.pg.GetAllRoutes(ctx)
	if err != nil {
		hs.logger.Error("Error get all routes", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(routes))
	ctx.Abort()
}

func (hs *handlerService) NewReviewRoute(ctx *gin.Context) {
	var paramsURI struct {
		RouteID string `uri:"routeId" binding:"required,uuid"`
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

	routeID, _ := uuid.Parse(paramsURI.RouteID)
	reviewRoute, err := hs.pg.NewReviewRoute(ctx, models.ReviewRoute{
		OwnerID:    user.ID,
		RouteID:    routeID,
		ReviewText: params.ReviewText,
		Stars:      params.Stars,
		CreatedAt:  time.Now(),
	})

	if err != nil {
		if hs.pg.IsError(pgerrcode.IsIntegrityConstraintViolation, err) {
			ctx.JSON(http.StatusConflict, models.NewErrorResponse(errs.NewConflict("You have already added a review to this route")))
		} else {
			hs.logger.Error("Error new review event", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		}
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(reviewRoute))
	ctx.Abort()
}

func (hs *handlerService) EditReviewRoute(ctx *gin.Context) {
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

func (hs *handlerService) GetReviewsRoutes(ctx *gin.Context) {
	var params struct {
		RouteID string `uri:"routeId" binding:"required,uuid"`
	}

	if response, statusCode, err := hs.validateAndShouldBindURI(ctx, &params); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}

	routeID, _ := uuid.Parse(params.RouteID)
	reviews, err := hs.pg.GetReviewsEvent(ctx, routeID)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		hs.logger.Debug("Error get reviews", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(reviews))
	ctx.Abort()
}
