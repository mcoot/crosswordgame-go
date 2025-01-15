package main

import (
	"context"
	"github.com/getkin/kin-openapi/openapi3"
	internalapi "github.com/mcoot/crosswordgame-go/internal/api"
	"github.com/mcoot/crosswordgame-go/internal/game"
	"github.com/mcoot/crosswordgame-go/internal/game/store"
	"github.com/mcoot/crosswordgame-go/internal/logging"
	"github.com/mcoot/crosswordgame-go/internal/utils"
	"log"
	"net/http"
	"slices"

	"github.com/oapi-codegen/nethttp-middleware"
)

func main() {
	ctx := utils.RootContext()
	ctx, err := logging.AddLoggerToContext(ctx, true)
	if err != nil {
		log.Fatalf("error adding logger to utils: %v", err)
	}
	logger := logging.GetLogger(ctx, "main")

	mux := http.NewServeMux()

	gameStore := store.NewInMemoryStore()
	gameManager := game.NewGameManager(gameStore)
	api := internalapi.NewCrosswordGameAPI(gameManager)
	api.AttachToMux(mux)

	h, err := setupMiddleware(ctx, mux)
	if err != nil {
		logger.Fatalf("error setting up middleware: %v", err)
	}

	logger.Infow("starting server", "port", 8080)
	if err := http.ListenAndServe(":8080", h); err != nil {
		logger.Fatalf("error serving: %v", err)
	}
}

func setupMiddleware(ctx context.Context, h http.Handler) (http.Handler, error) {
	openApiMiddleware, err := buildOpenApiMiddleware(ctx, "./schema/openapi.yaml")
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
