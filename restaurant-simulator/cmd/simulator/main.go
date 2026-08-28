// Command simulator is a standalone stub restaurant integration: it syncs a
// fixed menu to the Avito.Kitchen main service, receives pushed orders over
// a webhook (with a polling fallback), and auto-accepts and progresses them
// through the delivery pipeline — simulating how a real third-party
// restaurant system would integrate with the platform.
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"avito-kitchen-restaurant-simulator/internal/client"
	"avito-kitchen-restaurant-simulator/internal/config"
	"avito-kitchen-restaurant-simulator/internal/httpapi"
	"avito-kitchen-restaurant-simulator/internal/logger"
	"avito-kitchen-restaurant-simulator/internal/simulate"
	"avito-kitchen-restaurant-simulator/internal/store"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	log, err := logger.New(cfg.LogLevel)
	if err != nil {
		return err
	}
	defer func() { _ = log.Sync() }()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	c := client.New(cfg.MainServiceURL, cfg.RestaurantID, 5*time.Second)
	st := store.New()
	proc := simulate.NewProcessor(c, st, log)

	go simulate.SyncMenuLoop(ctx, c, cfg.MenuSyncInterval, log)
	go simulate.PollLoop(ctx, c, proc, cfg.PollInterval, log)

	httpServer := &http.Server{
		Addr:              ":" + cfg.ListenPort,
		Handler:           httpapi.NewRouter(proc, log),
		ReadHeaderTimeout: 5 * time.Second,
	}

	serveErr := make(chan error, 1)
	go func() {
		log.Info("webhook server listening", zap.String("addr", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serveErr <- err
			return
		}
		serveErr <- nil
	}()

	select {
	case <-ctx.Done():
		log.Info("shutdown signal received")
	case err := <-serveErr:
		if err != nil {
			return fmt.Errorf("http server: %w", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return httpServer.Shutdown(shutdownCtx)
}
