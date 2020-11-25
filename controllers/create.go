package controllers

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/steevehook/expenses-rest-api/logging"
	"github.com/steevehook/expenses-rest-api/models"
	"github.com/steevehook/expenses-rest-api/transport"
)

type expenseCreator interface {
	CreateExpense(models.CreateExpenseRequest) error
}

func createExpense(service expenseCreator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req models.CreateExpenseRequest
		err := parseBody(r, &req)
		if err != nil {
			logging.Logger.Error("could not unmarshal create expense body", zap.Error(err))
			transport.SendHTTPError(w, err)
			return
		}
		if err := req.Validate(); err != nil {
			transport.SendHTTPError(w, err)
			return
		}

		err = service.CreateExpense(req)
		if err != nil {
			logging.Logger.Debug("could not create expense", zap.Error(err))
			transport.SendHTTPError(w, err)
			return
		}

		logging.Logger.Info("successfully created expense")
		w.WriteHeader(http.StatusCreated)
	})
}
