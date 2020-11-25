package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/steevehook/expenses-rest-api/models"
	"github.com/steevehook/expenses-rest-api/transport"
)

const (
	pageQueryParam     = "page"
	pageSizeQueryParam = "page_size"

	defaultPage     = 1
	defaultPageSize = 10
)

type allExpensesGetter interface {
	GetAllExpenses(models.GetAllExpensesRequest) ([]models.Expense, error)
	ExpensesCount() (int, error)
}

type getAllExpensesResponse struct {
	Items    []models.Expense `json:"items"`
	Total    int              `json:"total"`
	NextPage string           `json:"next_page,omitempty"`
	PrevPage string           `json:"prev_page,omitempty"`
}

func getAllExpenses(service allExpensesGetter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page, err := parseQueryParam(r, pageQueryParam, defaultPage)
		if err != nil {
			transport.SendHTTPError(w, err)
			return
		}
		pageSize, err := parseQueryParam(r, pageSizeQueryParam, defaultPageSize)
		if err != nil {
			transport.SendHTTPError(w, err)
			return
		}
		req := models.GetAllExpensesRequest{
			Page:     page,
			PageSize: pageSize,
		}

		expenses, err := service.GetAllExpenses(req)
		if err != nil {
			transport.SendHTTPError(w, err)
			return
		}
		count, err := service.ExpensesCount()
		if err != nil {
			transport.SendHTTPError(w, err)
			return
		}
		res := getAllExpensesResponse{
			Items: expenses,
			Total: count,
		}
		if page*pageSize+1 <= count {
			res.NextPage = fmt.Sprintf("%s?%s", r.URL.Path, encodePageParams(page+1, pageSize))
		}
		if (page-1)*pageSize < count && page-1 > 0 {
			res.PrevPage = fmt.Sprintf("%s?%s", r.URL.Path, encodePageParams(page-1, pageSize))
		}
		transport.SendJSON(w, http.StatusOK, res)
	})
}

func parseQueryParam(r *http.Request, paramName string, defaultValue int) (int, error) {
	param := r.URL.Query().Get(paramName)
	if strings.TrimSpace(param) == "" {
		return defaultValue, nil
	}
	intParam, err := strconv.Atoi(param)
	if err != nil || intParam < 1 {
		e := models.DataValidationError{
			Message: fmt.Sprintf("invalid value: %s for param: %s", param, paramName),
		}
		return 0, e
	}
	return intParam, nil
}

func encodePageParams(page, pageSize int) string {
	params := url.Values{}
	params.Add(pageQueryParam, strconv.Itoa(page))
	params.Add(pageSizeQueryParam, strconv.Itoa(pageSize))
	return params.Encode()
}
