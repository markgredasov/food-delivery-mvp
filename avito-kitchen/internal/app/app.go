package app

import (
	"context"
	"fmt"
	"os"

	database_postgres "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/database/postgres"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/logger"
	server_http "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/server/http/server"
	"go.uber.org/zap"
)

type App struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	logger     *logger.Logger
	server     *server_http.HTTPServer
	pool       *database_postgres.ConnectionPool
}

func NewApp(ctx context.Context, cancelFunc context.CancelFunc) *App {
	return &App{
		ctx:        ctx,
		cancelFunc: cancelFunc,
	}
}

func (a *App) Run() error {
	if err := a.initLogger(); err != nil {
		//nolint:forbidigo // logger was not initialized to write here smth
		fmt.Println("failed to initialize logger:", err)
		os.Exit(1)
	}

	a.logger.Debug("initializing connection pool")
	if err := a.initPostgresPool(); err != nil {
		a.logger.Fatal("initialize pool", zap.Error(err))
	}

	a.logger.Debug("initializing features")
	routes := a.initFeatures()

	a.logger.Debug("initializing HTTP server")
	a.initHTTPServer(routes)

	return a.server.Run(a.ctx)
}

func (a *App) Close() {
	a.pool.Close()
	a.logger.Close()
}
