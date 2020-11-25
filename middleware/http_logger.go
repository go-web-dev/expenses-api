package middleware

import (
	"net/http"

	"github.com/steevehook/expenses-rest-api/logging"
	"go.uber.org/zap"
)

// HTTPLogger logs http requests
func HTTPLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logging.Logger.Debug(
			"http request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("query", r.URL.Query().Encode()),
		)
		h.ServeHTTP(w, r)
	})
}
