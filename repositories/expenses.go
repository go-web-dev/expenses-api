package repositories

import (
	"github.com/steevehook/expenses-rest-api/models"
)

// Closer represents application db closer
type Closer interface {
	Close() error
}

type Expenses interface {
	GetAllExpenses(page, size int) ([]models.Expense, error)
	GetExpensesByIDs(ids []string) ([]models.Expense, error)
	CreateExpense(title, currency string, price float64) error
	UpdateExpense(id, title, currency string, price float64) error
	DeleteExpense(id string) error
	Count() (int, error)
	Closer
}
