package middleware

import (
	"net/http"

	"github.com/steevehook/expenses-rest-api/logging"
	"github.com/steevehook/expenses-rest-api/models"
	"github.com/steevehook/expenses-rest-api/transport"
)

// JSONBody rejects endpoints that have missing application/json content type
func JSONBody(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(models.ContentType) != models.ApplicationJSONType {
			err := models.InvalidJSONError{
				Message: "missing json body or json content-type",
			}
			logging.Logger.Error("missing " + models.ApplicationJSONType + " content type")
			transport.SendHTTPError(w, err)
			return
		}
		h.ServeHTTP(w, r)
	})
}
