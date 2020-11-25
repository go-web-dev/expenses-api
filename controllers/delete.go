package controllers

import (
	"net/http"

	"github.com/steevehook/expenses-rest-api/logging"
	"github.com/steevehook/expenses-rest-api/transport"
)

type expenseDeleter interface {
	DeleteExpense(id string) error
}

func deleteExpense(service expenseDeleter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := parseIDParam(r)
		if err != nil {
			transport.SendHTTPError(w, err)
			return
		}

		err = service.DeleteExpense(id)
		if err != nil {
			transport.SendHTTPError(w, err)
			return
		}
		logging.Logger.Info("successfully deleted expense")
		w.WriteHeader(http.StatusNoContent)
	})
}
