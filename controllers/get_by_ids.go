package controllers

import (
	"net/http"

	"github.com/steevehook/expenses-rest-api/models"
	"github.com/steevehook/expenses-rest-api/transport"
)

type expensesByIDsGetter interface {
	GetExpensesByIDs(request models.GetExpensesByIDsRequest) ([]models.Expense, error)
}

type getExpensesByIDsResponse struct {
	Items []models.Expense `json:"items"`
}

func getExpensesByIDs(service expensesByIDsGetter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ids, err := parseIDsParam(r)
		if err != nil {
			transport.SendHTTPError(w, err)
			return
		}

		req := models.GetExpensesByIDsRequest{
			IDs: ids,
		}
		expenses, err := service.GetExpensesByIDs(req)
		if err != nil {
			transport.SendHTTPError(w, err)
			return
		}

		res := getExpensesByIDsResponse{
			Items: expenses,
		}
		transport.SendJSON(w, http.StatusOK, res)
	})
}
