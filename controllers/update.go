package controllers

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/steevehook/expenses-rest-api/logging"
	"github.com/steevehook/expenses-rest-api/models"
	"github.com/steevehook/expenses-rest-api/transport"
)

type expenseUpdater interface {
	UpdateExpense(models.UpdateExpenseRequest) error
}

func updateExpense(service expenseUpdater) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req models.UpdateExpenseRequest
		err := parseBody(r, &req)
		if err != nil {
			logging.Logger.Error("could not unmarshal update expense body", zap.Error(err))
			transport.SendHTTPError(w, err)
			return
		}
		if err := req.Validate(); err != nil {
			transport.SendHTTPError(w, err)
			return
		}

		id, err := parseIDParam(r)
		if err != nil {
			transport.SendHTTPError(w, err)
			return
		}
		req.ID = id

		err = service.UpdateExpense(req)
		if err != nil {
			transport.SendHTTPError(w, err)
			return
		}
		logging.Logger.Info("successfully updated the expense")
		w.WriteHeader(http.StatusNoContent)
	})
}
