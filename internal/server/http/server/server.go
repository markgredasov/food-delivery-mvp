package server_http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/markgredasov/food-delivery-mvp/internal/logger"
	"github.com/markgredasov/food-delivery-mvp/internal/server/http/middleware"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux        *http.ServeMux
	config     Config
	log        *logger.Logger
	middleware []middleware.Middleware
}

func NewHTTPServer(config Config, logger *logger.Logger, middleware ...middleware.Middleware) *HTTPServer {
	return &HTTPServer{
		mux:        http.NewServeMux(),
		config:     config,
		log:        logger,
		middleware: middleware,
	}
}

func (h *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)

		h.mux.Handle(prefix+"/", http.StripPrefix(prefix, router))
	}
}

func (h *HTTPServer) Run(ctx context.Context) error {
	mux := middleware.ChainMiddleware(h.mux, h.middleware...)

	server := &http.Server{
		Addr:    h.config.Address,
		Handler: mux,
	}

	ch := make(chan error, 1)

	go func() {
		defer func() {
			close(ch)
		}()

		h.log.Warn(
			"HTTP server started",
			zap.String("address", h.config.Address),
		)

		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		return fmt.Errorf("listen and serve HTTP server: %w", err)
	case <-ctx.Done():
		h.log.Warn("shutdown HTTP server")
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			h.config.ShutdownTimeout,
		)

		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		h.log.Warn("HTTP server stopped")
	}

	return nil
}
