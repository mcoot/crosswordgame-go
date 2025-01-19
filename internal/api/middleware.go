package api

import (
	"context"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/mcoot/crosswordgame-go/internal/logging"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
	"net/http"
	"slices"
)

func setupMiddleware(ctx context.Context, h http.Handler, schemaPath string) (http.Handler, error) {
	openApiMiddleware, err := buildOpenApiMiddleware(ctx, schemaPath)
	if err != nil {
		return nil, err
	}

	// Middleware list in first-to-last order
	middlewares := []func(next http.Handler) http.Handler{
		ctxMiddleware(ctx),
		requestLoggerMiddleware,
		openApiMiddleware,
	}

	slices.Reverse(middlewares)
	for _, m := range middlewares {
		h = m(h)
	}

	return h, nil
}

func ctxMiddleware(ctx context.Context) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func buildOpenApiMiddleware(ctx context.Context, schemaPath string) (func(next http.Handler) http.Handler, error) {
	loader := openapi3.Loader{
		Context: ctx,
	}
	doc, err := loader.LoadFromFile(schemaPath)
	if err != nil {
		return nil, err
	}
	err = doc.Validate(ctx)
	if err != nil {
		return nil, err
	}

	return nethttpmiddleware.OapiRequestValidatorWithOptions(
		doc,
		&nethttpmiddleware.Options{
			ErrorHandler: func(w http.ResponseWriter, message string, statusCode int) {
				logger := logging.GetLogger(ctx, "openapi")
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
		logger := logging.GetLogger(r.Context(), "request")
		logger.Infow("request", "method", r.Method, "url", r.URL)
		next.ServeHTTP(w, r)
	})
}
