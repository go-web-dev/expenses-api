package models

import (
	"strings"
)

// GetAllExpensesRequest represents http request for fetching all expenses with pagination
type GetAllExpensesRequest struct {
	Page     int
	PageSize int
}

// GetExpensesByIDsRequest represents http request for fetching a list of expenses by ids
type GetExpensesByIDsRequest struct {
	IDs []string
}

// CreateExpenseRequest represents http request for creating an expense
type CreateExpenseRequest struct {
	Title    string  `json:"title"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
}

func (r CreateExpenseRequest) Validate() error {
	return validateExpenseReqBody(r.Title, r.Currency, r.Price, false)
}

// UpdateExpenseRequest represents http request for updating an expense
type UpdateExpenseRequest struct {
	ID       string
	Title    string  `json:"title"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
}

func (r UpdateExpenseRequest) Validate() error {
	return validateExpenseReqBody(r.Title, r.Currency, r.Price, true)
}

func validateExpenseReqBody(title, currency string, price float64, optional bool) error {
	if !optional && strings.TrimSpace(title) == "" {
		return DataValidationError{Message: "title should not be empty"}
	}

	if !optional && price <= 0 || optional && price < 0 {
		return DataValidationError{Message: "price must be greater than 0"}
	}

	currencies := []string{"USD", "EUR", "GBP", "MDL"}
	c := strings.TrimSpace(strings.ToUpper(currency))
	if len(c) == 0 && optional {
		return nil
	}
	switch c {
	case currencies[0], currencies[1], currencies[2], currencies[3]:
		return nil
	default:
		return DataValidationError{Message: "currency must be one of: " + strings.Join(currencies, ",")}
	}
}
