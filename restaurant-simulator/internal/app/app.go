package app

import (
	"avito-kitchen-restaurant-simulator/internal/config"
	"avito-kitchen-restaurant-simulator/internal/logger"
	server_http "avito-kitchen-restaurant-simulator/internal/server/http/server"
	"context"
	"fmt"
	"os"
)

type App struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	logger     *logger.Logger
	server     *server_http.HTTPServer
	cfg        config.Config
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

	a.logger.Debug("initializing config")
	cfg, err := config.Load()
	if err != nil {
		a.logger.Fatal("cannot initialize config")
	}
	a.cfg = cfg

	a.logger.Debug("initializing features")
	routes := a.initFeatures()

	a.logger.Debug("initializing HTTP server")
	a.initHTTPServer(routes)

	return a.server.Run(a.ctx)
}

func (a *App) Close() {
	a.logger.Close()
}
