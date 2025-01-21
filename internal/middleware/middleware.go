package middleware

import (
	"context"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/mux"
	"github.com/mcoot/crosswordgame-go/internal/logging"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
	"go.uber.org/zap"
	"net/http"
	"slices"
)

func SetupRoutedMiddleware(router *mux.Router, baseLogger *zap.SugaredLogger, schemaPath string) error {
	openApiMiddleware, err := buildOpenApiMiddleware(
		baseLogger.Named("openapi"),
		schemaPath,
	)
	if err != nil {
		return err
	}

	router.Use(openApiMiddleware)

	return nil
}

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

func buildOpenApiMiddleware(logger *zap.SugaredLogger, schemaPath string) (func(next http.Handler) http.Handler, error) {
	loader := openapi3.Loader{
		Context: context.Background(),
	}
	doc, err := loader.LoadFromFile(schemaPath)
	if err != nil {
		return nil, err
	}
	err = doc.Validate(context.Background())
	if err != nil {
		return nil, err
	}

	return nethttpmiddleware.OapiRequestValidatorWithOptions(
		doc,
		&nethttpmiddleware.Options{
			ErrorHandler: func(w http.ResponseWriter, message string, statusCode int) {
				logger.Warnw(
					"openapi validation error",
					"message", message,
					"status_code", statusCode,
				)
				http.Error(w, message, statusCode)
			},
		},
	), nil
}

func requestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logging.GetLogger(r.Context()).Named("access_log")
		logger.Infow("access log", "method", r.Method, "url", r.URL)
		next.ServeHTTP(w, r)
	})
}
