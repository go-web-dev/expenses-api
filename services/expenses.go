package services

import (
	"go.uber.org/zap"

	"github.com/steevehook/expenses-rest-api/logging"
	"github.com/steevehook/expenses-rest-api/models"
	"github.com/steevehook/expenses-rest-api/repositories"
)

// Expenses represents the Expenses service
type Expenses struct {
	ExpensesRepo repositories.Expenses
}

// GetAllExpenses fetches all expenses with pagination possibilities
func (s Expenses) GetAllExpenses(req models.GetAllExpensesRequest) ([]models.Expense, error) {
	expenses, err := s.ExpensesRepo.GetAllExpenses(req.Page, req.PageSize)
	if err != nil {
		logging.Logger.Error("could not fetch all expenses from db", zap.Error(err))
		return []models.Expense{}, err
	}
	return expenses, nil
}

// GetExpensesByIDs fetches expenses by a list of given IDs
func (s Expenses) GetExpensesByIDs(req models.GetExpensesByIDsRequest) ([]models.Expense, error) {
	expenses, err := s.ExpensesRepo.GetExpensesByIDs(req.IDs)
	if err != nil {
		logging.Logger.Error("could not fetch expenses by ids from db", zap.Error(err))
		return []models.Expense{}, err
	}
	return expenses, nil
}

// CreateExpense creates a brand new expense
func (s Expenses) CreateExpense(req models.CreateExpenseRequest) error {
	err := s.ExpensesRepo.CreateExpense(req.Title, req.Currency, req.Price)
	if err != nil {
		logging.Logger.Error("could not create expense in db", zap.Error(err))
	}
	return nil
}

// UpdateExpense updates an existing created expense
func (s Expenses) UpdateExpense(req models.UpdateExpenseRequest) error {
	err := s.ExpensesRepo.UpdateExpense(req.ID, req.Title, req.Currency, req.Price)
	if err != nil {
		logging.Logger.Error("could not update expense in db", zap.Error(err))
	}
	return nil
}

// DeleteExpense deletes a list of expenses by a given list of IDs
func (s Expenses) DeleteExpense(id string) error {
	return s.ExpensesRepo.DeleteExpense(id)
}

// ExpensesCount fetches the total count of created expenses
func (s Expenses) ExpensesCount() (int, error) {
	return s.ExpensesRepo.Count()
}
