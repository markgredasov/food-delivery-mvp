package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/logger"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/server/http/response"
	"go.uber.org/zap"
)

var (
	loggerKey string = "logger"
)

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		const requestIDHeader = "x-request-id"
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set(requestIDHeader, requestID)
			w.Header().Set(requestIDHeader, requestID)

			next.ServeHTTP(w, r)
		})
	}
}

func Logger(l *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		const requestIDHeader = "x-request-id"
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)

			logger := l.With(
				zap.String("Request-ID", requestID),
				zap.String("Method", r.Method),
				zap.String("URL", r.URL.String()),
			)

			//nolint:staticcheck // not needed
			ctx := context.WithValue(r.Context(), loggerKey, logger)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Recovery() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logger := logger.FromContext(ctx)
			responseHandler := response.NewHTTPResponseHandler(logger, w)

			defer func() {
				if r := recover(); r != nil {
					responseHandler.PanicResponse(
						r,
						"during handle HTTP request got unexpected panic",
					)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logger := logger.FromContext(ctx)

			writer := response.NewResponseWriter(w)

			start := time.Now().UTC()
			logger.Debug(
				">>> incoming HTTP request",
				zap.Time("time", start),
			)

			next.ServeHTTP(writer, r)

			logger.Debug(
				"<<< done HTTP request",
				zap.Int("status-code", writer.GetStatusCode()),
				zap.Duration("latency", time.Since(start)),
			)
		})
	}
}

func CORS() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowedOrigins := map[string]struct{}{
				"http://localhost:8080": {},
			}

			origin := r.Header.Get("Origin")

			if _, ok := allowedOrigins[origin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
