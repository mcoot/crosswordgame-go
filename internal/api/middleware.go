package api

import (
	"github.com/mcoot/crosswordgame-go/internal/logging"
	"go.uber.org/zap"
	"net/http"
	"slices"
)

func SetupGlobalMiddleware(h http.Handler, baseLogger *zap.SugaredLogger) http.Handler {
	middlewares := []func(next http.Handler) http.Handler{
		loggerInContextMiddleware(baseLogger),
		requestLoggerMiddleware,
	}
	slices.Reverse(middlewares)

	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func loggerInContextMiddleware(baseLogger *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := baseLogger.
				Named("api").
				With("path", r.URL.Path)
			reqCtx := r.Context()

			ctxWithLogger := logging.AddLoggerToContext(reqCtx, logger)

			r = r.WithContext(ctxWithLogger)
			next.ServeHTTP(w, r)
		})
	}
}

func requestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logging.GetLogger(r.Context()).Named("access_log")
		logger.Infow("access log", "method", r.Method, "url", r.URL)
		next.ServeHTTP(w, r)
	})
}
