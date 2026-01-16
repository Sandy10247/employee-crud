package middleware

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type contextLoggerKey string

const (
	// UserIDKey is the key for user ID in the request context
	loggerKey contextLoggerKey = "logger"
)

func ZapMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			// Add the logger to the request context
			ctx := context.WithValue(r.Context(), loggerKey, logger)
			r = r.WithContext(ctx)

			// Log the request details after the handler finishes
			defer func() {
				logger.Info("request completed",
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
					zap.Duration("duration", time.Since(start)),
				)
			}()

			next.ServeHTTP(w, r)
		})
	}
}
