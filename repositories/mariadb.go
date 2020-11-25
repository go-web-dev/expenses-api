package repositories

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mysql"
	"go.uber.org/zap"

	"github.com/steevehook/expenses-rest-api/logging"
	"github.com/steevehook/expenses-rest-api/models"
)

const (
	expensesTableName = "expenses"
)

// MariaDBSettings represents the settings for MariaDB
type MariaDBSettings struct {
	URL                string
	MaxOpenConnections int
	MaxIdleConnections int
	ConnMaxLifetime    time.Duration
}

// MariaDBDriver represents MariaDB repository driver
type MariaDBDriver struct {
	mariaDB db.Session
}

// NewMariaDBDriver creates a new instance of MariaDB database
func NewMariaDBDriver(settings MariaDBSettings) (*MariaDBDriver, error) {
	conn, err := mysql.ParseURL(settings.URL)
	if err != nil {
		logging.Logger.Error("could not parse mariadb connection url", zap.Error(err))
		return nil, err
	}
	session, err := mysql.Open(conn)
	if err != nil {
		logging.Logger.Error("could not open mariadb database", zap.Error(err))
		return nil, err
	}
	session.SetConnMaxLifetime(settings.ConnMaxLifetime)
	session.SetMaxOpenConns(settings.MaxOpenConnections)
	session.SetMaxIdleConns(settings.MaxIdleConnections)
	driver := &MariaDBDriver{
		mariaDB: session,
	}
	return driver, nil
}

// GetAllExpenses fetches all expenses with pagination possibilities from MariaDB
func (d MariaDBDriver) GetAllExpenses(page, pageSize int) ([]models.Expense, error) {
	var expenses []models.Expense
	err := d.mariaDB.
		Collection(expensesTableName).
		Find().
		Page(uint(page)).
		Paginate(uint(pageSize)).
		OrderBy("modified_at").
		All(&expenses)
	if err != nil {
		logging.Logger.Error("could not execute find all on mariadb expenses records", zap.Error(err))
		return []models.Expense{}, err
	}
	return expenses, nil
}

// GetExpensesByIDs fetches a list of expenses by a given list of IDs from MariaDB
func (d MariaDBDriver) GetExpensesByIDs(ids []string) ([]models.Expense, error) {
	var expenses []models.Expense
	idsPlaceholder := strings.Repeat("?,", len(ids)-1)
	idsPlaceholder += "?"
	var args []interface{}
	args = append(args, fmt.Sprintf("id IN(%s)", idsPlaceholder))
	for _, id := range ids {
		args = append(args, id)
	}
	err := d.mariaDB.
		SQL().
		SelectFrom(expensesTableName).
		Where(args...).
		All(&expenses)
	if err != nil {
		logging.Logger.Error("could not select expense records from mariadb", zap.Error(err))
		return []models.Expense{}, err
	}
	return expenses, nil
}

// CreateExpense creates a brand new expense and saves it into MariaDB
func (d MariaDBDriver) CreateExpense(title, currency string, price float64) error {
	uid := uuid.New()
	expense := models.Expense{
		ID:         uid.String(),
		Title:      title,
		Currency:   currency,
		Price:      price,
		CreatedAt:  time.Now().UTC(),
		ModifiedAt: time.Now().UTC(),
	}
	_, err := d.mariaDB.Collection(expensesTableName).Insert(expense)
	if err != nil {
		logging.Logger.Error("could not create expense record in mariadb", zap.Error(err))
		return err
	}
	return nil
}

// UpdateExpense updates an existing expense and updates the record in MariaDB
func (d MariaDBDriver) UpdateExpense(id, title, currency string, price float64) error {
	var modified bool
	expense, err := d.findExpense(id)
	if err != nil {
		return err
	}
	if title != "" && expense.Title != title {
		expense.Title = title
		modified = true
	}
	if currency != "" && expense.Currency != currency {
		expense.Currency = currency
		modified = true
	}
	if price > 0 && expense.Price != price {
		expense.Price = price
		modified = true
	}
	if modified {
		expense.ModifiedAt = time.Now().UTC()
	}
	err = d.mariaDB.Collection(expensesTableName).UpdateReturning(&expense)
	if err != nil {
		logging.Logger.Error("could not update expense in mariadb", zap.Error(err))
		return err
	}
	return nil
}

// DeleteExpense deletes a given expense from MariaDB given a list of IDs
func (d MariaDBDriver) DeleteExpense(id string) error {
	_, err := d.findExpense(id)
	if err != nil {
		return err
	}
	_, err = d.mariaDB.
		SQL().
		DeleteFrom(expensesTableName).
		Where(db.Cond{"id": id}).
		Exec()
	if err != nil {
		logging.Logger.Error("could not delete expense from mariadb", zap.Error(err))
		return err
	}
	return nil
}

// Count fetches the total count from expenses table from MariaDB
func (d MariaDBDriver) Count() (int, error) {
	count, err := d.mariaDB.Collection(expensesTableName).Count()
	return int(count), err
}

// Close closes the MariaDB database
func (d MariaDBDriver) Close() error {
	logging.Logger.Info("stopping mariadb server")
	err := d.mariaDB.Close()
	if err != nil {
		return err
	}

	logging.Logger.Info("mariadb server successfully stopped")
	return nil
}

func (d MariaDBDriver) findExpense(id string) (models.Expense, error) {
	var expense models.Expense
	err := d.mariaDB.Collection(expensesTableName).
		Find(db.Cond{"id": id}).
		One(&expense)
	if err != nil {
		logging.Logger.Debug("could not find expense in mariadb", zap.String("id", id))
		e := models.ResourceNotFoundError{
			Message: fmt.Sprintf("could not find expense with id: %s", id),
		}
		return models.Expense{}, e
	}
	return expense, nil
}
