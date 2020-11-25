package transport

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/steevehook/expenses-rest-api/logging"
	"github.com/steevehook/expenses-rest-api/models"
)

const (
	// DataValidationErrorType describes data validation errors
	DataValidationErrorType = "data_validation_error"
	// FormatValidationErrorType describes format validation errors
	FormatValidationErrorType = "format_validation_error"
	// ResourceNotFoundErrorType describes a severe resource not found
	ResourceNotFoundErrorType = "resource_not_found"
	// ServiceErrorType describes a severe generic server error
	ServiceErrorType = "service_error"
)

// SendHTTPError converts errors into HTTP JSON errors
func SendHTTPError(w http.ResponseWriter, err error) {
	httpError := toHTTPError(err)
	SendJSON(w, httpError.Code, httpError)
}

// SendJSON converts application response into JSON responses
func SendJSON(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set(models.ContentType, models.ApplicationJSONType)
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(&response); err != nil {
		logging.Logger.Error("could not encode response", zap.Error(err))
	}
}

func toHTTPError(err error) models.HTTPError {
	switch e := err.(type) {
	case models.HTTPError:
		return e

	case models.InvalidJSONError:
		return models.HTTPError{
			Code:    http.StatusBadRequest,
			Type:    FormatValidationErrorType,
			Message: e.Message,
		}

	case models.FormatValidationError:
		return models.HTTPError{
			Code:    http.StatusBadRequest,
			Type:    FormatValidationErrorType,
			Message: e.Message,
		}

	case models.DataValidationError:
		return models.HTTPError{
			Code:    http.StatusBadRequest,
			Type:    DataValidationErrorType,
			Message: e.Message,
		}

	case models.ResourceNotFoundError:
		return models.HTTPError{
			Code:    http.StatusNotFound,
			Type:    ResourceNotFoundErrorType,
			Message: e.Message,
		}

	default:
		return models.HTTPError{
			Code:    http.StatusInternalServerError,
			Type:    ServiceErrorType,
			Message: "server was not able to process your request",
		}
	}
}
