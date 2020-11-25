package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"go.uber.org/zap"

	"github.com/steevehook/expenses-rest-api/config"
	"github.com/steevehook/expenses-rest-api/controllers"
	"github.com/steevehook/expenses-rest-api/logging"
	"github.com/steevehook/expenses-rest-api/models"
	"github.com/steevehook/expenses-rest-api/repositories"
	"github.com/steevehook/expenses-rest-api/services"
)

type App struct {
	stopOnce sync.Once
	Server   *http.Server
	Cfg      *config.Manager
	dbCloser repositories.Closer
}

// Init initializes the application
func Init(configPath string) (*App, error) {
	configManager, err := config.Init(configPath)
	if err != nil {
		return nil, fmt.Errorf("could not initialize app config: %v", err)
	}
	if err := logging.Init(configManager); err != nil {
		return nil, fmt.Errorf("could not initialize logger: %v", err)
	}

	var expensesRepo repositories.Expenses
	switch configManager.AppDBType() {
	case models.BoltDBType:
		expensesRepo, err = repositories.NewBoltDriver(configManager.BoltDBFileName())
		if err != nil {
			return nil, err
		}
	case models.MariaDBType:
		dbSettings := repositories.MariaDBSettings{
			URL:                configManager.MariaDBUrl(),
			MaxOpenConnections: configManager.MariaDBMaxOpenConnections(),
			MaxIdleConnections: configManager.MariaDBMaxIdleConnections(),
			ConnMaxLifetime:    configManager.MariaDBConnMaxLifetime(),
		}
		expensesRepo, err = repositories.NewMariaDBDriver(dbSettings)
		if err != nil {
			return nil, err
		}
	}

	routerCfg := controllers.RouterConfig{
		ExpensesSvc: services.Expenses{
			ExpensesRepo: expensesRepo,
		},
	}
	app := &App{
		Cfg: configManager,
		Server: &http.Server{
			Addr:         configManager.AppListen(),
			Handler:      controllers.NewRouter(routerCfg),
			ReadTimeout:  configManager.AppReadTimeout(),
			WriteTimeout: configManager.AppWriteTimeout(),
			ErrorLog:     logging.HTTPServerLogger(),
		},
		dbCloser: expensesRepo,
	}
	return app, nil
}

// Start starts the application
func (a *App) Start() error {
	logging.Logger.Info(
		"http server is ready to handle requests",
		zap.String("listen", a.Cfg.AppListen()),
	)

	err := a.Server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Stop shuts down the http server
func (a *App) Stop() error {
	var err error
	a.stopOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), a.Cfg.AppShutdownTimeout())
		defer cancel()

		logging.Logger.Info("shutting down the http server")
		if e := a.Server.Shutdown(ctx); err != nil {
			logging.Logger.Error("error on server shutdown", zap.Error(e))
			err = e
			return
		}

		logging.Logger.Info("http server was shut down")

		err = a.dbCloser.Close()
		if err != nil {
			logging.Logger.Error("could not stop db", zap.Error(err))
		}
	})
	return err
}

// Stopper represents app stop feature
type Stopper interface {
	Stop() error
}

// ListenToSignals listens for any incoming termination signals and shuts down the application
func ListenToSignals(signals []os.Signal, apps ...Stopper) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, signals...)

	<-c
	for _, a := range apps {
		err := a.Stop()
		if err != nil {
			logging.Logger.Error("stopping resulted in error", zap.Error(err))
		}
	}

	os.Exit(0)
}
