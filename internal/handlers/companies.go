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
)

func (hs *handlerService) NewCompany(ctx *gin.Context) {
	var params struct {
		Name        string `json:"name" binding:"required,min=6"`
		Description string `json:"description" binding:"required,min=12"`
		PhotoCard   string `json:"photo_card" binding:"required,url"`
	}

	if response, statusCode, err := hs.validateAndShouldBindJSON(ctx, &params); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}
	vkParams := hs.GetVKParams(ctx)

	user, err := hs.pg.GetUserByVkID(ctx, int64(vkParams.VkUserID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	company, err := hs.pg.NewCompany(ctx, models.Company{
		UserID:      user.ID,
		Name:        params.Name,
		Description: params.Description,
		PhotoCard:   params.PhotoCard,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(company))
	ctx.Abort()
}

func (hs *handlerService) AcceptCompany(ctx *gin.Context) {
	var params struct {
		CompanyID string `uri:"companyId" binding:"required,uuid"`
	}

	if response, statusCode, err := hs.validateAndShouldBindURI(ctx, &params); err != nil {
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

	if !user.IsAdmin {
		ctx.JSON(http.StatusForbidden, models.NewErrorResponse(errs.NewForbidden("You don't have access to this method")))
		ctx.Abort()

		return
	}

	companyID, _ := uuid.Parse(params.CompanyID)
	company, err := hs.pg.GetCompanyByID(ctx, companyID)

	if err != nil {
		hs.logger.Error("Error get company", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	company.IsReleased = true
	if err = hs.pg.SaveCompany(ctx, company); err != nil {
		hs.logger.Error("Error save company", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(company))
	ctx.Abort()
}

func (hs *handlerService) GetCompany(ctx *gin.Context) {
	var params struct {
		CompanyID string `uri:"companyId" binding:"required,uuid"`
	}

	if response, statusCode, err := hs.validateAndShouldBindURI(ctx, &params); err != nil {
		ctx.JSON(statusCode, response)
		ctx.Abort()

		return
	}

	companyID, _ := uuid.Parse(params.CompanyID)
	company, err := hs.pg.GetCompanyByID(ctx, companyID)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		hs.logger.Error("Error get company", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	if !company.IsNil() {
		companyRating, err := hs.pg.GetCompanyAverageRating(ctx, company.ID)
		if err != nil {
			hs.logger.Error("Error with get rating company", zap.Error(err))

			companyRating = 0
		}

		ctx.JSON(http.StatusOK, models.NewResponse(struct {
			*models.Company
			Rating float64 `json:"rating"`
		}{
			Company: company,
			Rating:  companyRating,
		}))
	} else {
		ctx.JSON(http.StatusOK, models.NewResponse(nil))
	}
	ctx.Abort()
}

func (hs *handlerService) GetMyCompanies(ctx *gin.Context) {
	vkParams := hs.GetVKParams(ctx)

	companies, err := hs.pg.GetCompaniesByVkID(ctx, int64(vkParams.VkUserID))
	if err != nil {
		hs.logger.Error("Error get companies by vk id", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(companies))
	ctx.Abort()
}

func (hs *handlerService) GetAllCompanies(ctx *gin.Context) {
	companies, err := hs.pg.GetAllCompanies(ctx)
	if err != nil {
		hs.logger.Error("Error get all companies", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(errs.NewInternalServer("Internal server error")))
		ctx.Abort()

		return
	}

	ctx.JSON(http.StatusOK, models.NewResponse(companies))
	ctx.Abort()
}
