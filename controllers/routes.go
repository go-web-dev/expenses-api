package controllers

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	"github.com/steevehook/expenses-rest-api/middleware"
)

const (
	idRouteParam  = "id"
	idsRouteParam = "ids"
)

// ExpensesService represents the Expenses service interface
type ExpensesService interface {
	allExpensesGetter
	expensesByIDsGetter
	expenseCreator
	expenseUpdater
	expenseDeleter
}

// AuthenticationService represents the Authentication service interface
type AuthenticationService interface {
	loginner
	signupper
	logoutter
}

// RouterConfig represents the application router config
type RouterConfig struct {
	ExpensesSvc ExpensesService
	AuthSvc     AuthenticationService
}

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
)

// NewRouter creates a new application HTTP router
func NewRouter(cfg RouterConfig) http.Handler {
	chain := alice.New(
		middleware.HTTPLogger,
	)
	jsonBodyChain := chain.Append(
		middleware.JSONBody,
	)
	route := func(h http.Handler) http.Handler {
		return chain.Then(h)
	}
	routeWithBody := func(h http.Handler) http.Handler {
		return jsonBodyChain.Then(h)
	}
	recordMetrics()

	router := httprouter.New()
	router.Handler(http.MethodGet, "/metrics", promhttp.Handler())
	router.Handler(http.MethodGet, "/expenses", route(getAllExpenses(cfg.ExpensesSvc)))
	router.Handler(http.MethodGet, "/expenses/:"+idsRouteParam, route(getExpensesByIDs(cfg.ExpensesSvc)))
	router.Handler(http.MethodPost, "/expenses", routeWithBody(createExpense(cfg.ExpensesSvc)))
	router.Handler(http.MethodPatch, "/expenses/:"+idRouteParam, routeWithBody(updateExpense(cfg.ExpensesSvc)))
	router.Handler(http.MethodDelete, "/expenses/:"+idRouteParam, route(deleteExpense(cfg.ExpensesSvc)))
	router.Handler(http.MethodPost, "/login", routeWithBody(login(cfg.AuthSvc)))
	router.Handler(http.MethodPost, "/signup", routeWithBody(signup(cfg.AuthSvc)))
	router.Handler(http.MethodPost, "/logout", routeWithBody(logout(cfg.AuthSvc)))
	router.NotFound = route(NotFound())

	return router
}
