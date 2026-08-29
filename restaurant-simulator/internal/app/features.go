package app

import (
	"avito-kitchen-restaurant-simulator/internal/client"
	"avito-kitchen-restaurant-simulator/internal/handler"
	server_http "avito-kitchen-restaurant-simulator/internal/server/http/server"
	"avito-kitchen-restaurant-simulator/internal/simulate"
	"avito-kitchen-restaurant-simulator/internal/store"
	"context"
	"time"
)

func (a *App) initFeatures() []server_http.Route {
	c := client.New(a.cfg.MainServiceURL, a.cfg.RestaurantID, 5*time.Second)
	st := store.New()
	proc := simulate.NewProcessor(c, st)
	handler := handler.New(proc)

	a.ctx = context.WithValue(a.ctx, "logger", a.logger)

	go simulate.SyncMenuLoop(a.ctx, c, a.cfg.MenuSyncInterval)
	go simulate.PollLoop(a.ctx, c, proc, a.cfg.PollInterval)

	return handler.Routes()
}
