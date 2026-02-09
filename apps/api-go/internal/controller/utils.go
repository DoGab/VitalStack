package controller

import (
	"errors"
	"net/http"

	"github.com/dogab/vitalstack/api/pkg/types"

	"github.com/danielgtaylor/huma/v2"
)

var (
	errDefaultInternalServerError = errors.New("something went wrong, please try again later")
)

func convertServiceErrorToHTTPError(err error) huma.StatusError {
	// Check if the error is a field validation error
	var fieldValidationError types.FieldValidationError
	if errors.As(err, &fieldValidationError) {
		errorDetail := &huma.ErrorDetail{
			Message:  fieldValidationError.GetMessage(),
			Location: fieldValidationError.GetLocation(),
			Value:    fieldValidationError.GetValue(),
		}
		return huma.Error422UnprocessableEntity("validation error", errorDetail)
	}
	// Check if the error is a not found error
	var serviceError types.ServiceError
	if errors.As(err, &serviceError) {
		switch serviceError.HTTPStatus() {
		case http.StatusNotFound:
			return huma.Error404NotFound(serviceError.Error())
		case http.StatusBadRequest:
			return huma.Error400BadRequest(serviceError.Error())
		case http.StatusUnauthorized:
			return huma.Error401Unauthorized(serviceError.Error())
		case http.StatusForbidden:
			return huma.Error403Forbidden(serviceError.Error())
		case http.StatusConflict:
			return huma.Error409Conflict(serviceError.Error())
		case http.StatusTooManyRequests:
			return huma.Error429TooManyRequests(serviceError.Error())
		default:
			return huma.Error500InternalServerError("internal server error", errDefaultInternalServerError)
		}
	}
	// Return 500 for any other errors
	return huma.Error500InternalServerError("internal server error", errDefaultInternalServerError)
}
