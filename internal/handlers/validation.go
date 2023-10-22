package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ShpullRequest/backend/internal/errs"
	"github.com/ShpullRequest/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
	"strconv"
)

func (hs *handlerService) validateAndShouldBindURI(ctx *gin.Context, obj any) (*models.ErrorResponse, int, error) {
	if err := ctx.ShouldBindUri(obj); err != nil {
		if errors.Is(err, strconv.ErrSyntax) {
			return models.NewErrorResponse(errs.NewBadRequest("Invalid type uri variable")), http.StatusBadRequest, err
		}

		return hs.parseShouldBindErrors(err)
	}

	return nil, 0, nil
}

func (hs *handlerService) validateAndShouldBindJSON(ctx *gin.Context, obj any) (*models.ErrorResponse, int, error) {
	if err := ctx.ShouldBindJSON(obj); err != nil {
		if errors.Is(err, io.EOF) {
			return models.NewErrorResponse(errs.NewBadRequest("Request body not provided")), http.StatusBadRequest, err
		}

		var jsonTypeError *json.UnmarshalTypeError
		if ok := errors.As(err, &jsonTypeError); ok {
			return models.NewErrorResponse(
					fmt.Sprintf("Field value \"%s\" must be %s", jsonTypeError.Field, jsonTypeError.Type),
				),
				http.StatusBadRequest, err
		}

		var jsonError *json.SyntaxError
		if ok := errors.As(err, &jsonError); ok {
			return models.NewErrorResponse(errs.NewBadRequest("Invalid json object")), http.StatusBadRequest, err
		}

		return hs.parseShouldBindErrors(err)
	}

	return nil, 0, nil
}

func (hs *handlerService) parseShouldBindErrors(err error) (*models.ErrorResponse, int, error) {
	if ok, errResponse := hs.parseValidationErrors(err); ok {
		return errResponse, http.StatusBadRequest, err
	}

	var sliceValidationErrors binding.SliceValidationError
	if ok := errors.As(err, &sliceValidationErrors); ok && len(sliceValidationErrors) > 0 {
		if ok, errResponse := hs.parseValidationErrors(sliceValidationErrors[0]); ok {
			return errResponse, http.StatusBadRequest, err
		}
	}

	return models.NewErrorResponse(errs.NewBadRequest(fmt.Sprint(err))), http.StatusBadRequest, err
}

func (hs *handlerService) parseValidationErrors(err error) (bool, *models.ErrorResponse) {
	var validationErrors validator.ValidationErrors
	if ok := errors.As(err, &validationErrors); ok && len(validationErrors) > 0 {
		fErr := validationErrors[0]

		var errResponse string
		if fErr.Param() == "" {
			errResponse = fErr.Tag()
		} else {
			errResponse = fmt.Sprintf("%s=%s", fErr.Tag(), fErr.Param())
		}

		return true, models.NewErrorResponse(
			errs.NewBadRequest(
				fmt.Sprintf("Field validation for \"%s\" failed on the '%s' tag.", fErr.Field(), errResponse),
			),
		)
	}

	return false, nil
}
