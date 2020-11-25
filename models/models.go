package models

import (
	"time"
)

// Expense represents the expense model
type Expense struct {
	ID         string    `json:"id" db:"id"`
	Price      float64   `json:"price" db:"price"`
	Title      string    `json:"title" db:"title"`
	Currency   string    `json:"currency" db:"currency"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	ModifiedAt time.Time `json:"modified_at" db:"modified_at"`
}
