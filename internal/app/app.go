package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/markgredasov/food-delivery-mvp/internal/app/closer"
	"github.com/markgredasov/food-delivery-mvp/internal/app/di"
	"github.com/markgredasov/food-delivery-mvp/internal/logger"
	"github.com/markgredasov/food-delivery-mvp/internal/server/http/middleware"
	server_http "github.com/markgredasov/food-delivery-mvp/internal/server/http/server"
)

type App struct {
	server      *server_http.HTTPServer
	diContainer *di.Container
}

func New(ctx context.Context) *App {
	a := &App{
		diContainer: di.New(),
	}

	a.initDeps(ctx)

	return a
}

func (a *App) initDeps(ctx context.Context) {
	inits := []func(context.Context){
		a.diContainer.InitPoolAndRunMigrations,
		a.initHTTPServer,
	}

	for _, fn := range inits {
		fn(ctx)
	}
}

func (a *App) Run(ctx context.Context) {
	if err := a.server.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Println("run application:", err)
	}

	closerCtx, closerCancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second) //nolint:mnd // todo
	defer closerCancel()

	if err := closer.CloseAll(closerCtx); err != nil {
		fmt.Println("failed to close deps:", err)
	}
}

func (a *App) initHTTPServer(ctx context.Context) {
	apiVersion1Router := server_http.NewAPIVersionRouter(server_http.APIVersion1)
	apiVersion1Router.RegisterRoutes(a.diContainer.Routes()...)

	log := logger.FromContext(ctx)

	serverConfig := server_http.NewConfigMust()
	server := server_http.NewHTTPServer(
		serverConfig,
		log,
		middleware.RequestID(),
		middleware.Logger(log),
		middleware.Recovery(),
		middleware.Trace(),
		middleware.CORS(),
	)

	server.RegisterAPIRouters(apiVersion1Router)
	server.RegisterSwagger()
	a.server = server
}
