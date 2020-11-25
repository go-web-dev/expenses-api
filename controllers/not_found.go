package controllers

import (
	"net/http"

	"github.com/steevehook/expenses-rest-api/models"
	"github.com/steevehook/expenses-rest-api/transport"
)

// NotFound represents the resource not found handler
func NotFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := models.ResourceNotFoundError{
			Message: "route not found",
		}
		transport.SendHTTPError(w, err)
	})
}
