package jsonapi

import (
	"context"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/mux"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
	"go.uber.org/zap"
	"net/http"
)

func setupMiddleware(router *mux.Router, baseLogger *zap.SugaredLogger, schemaPath string) error {
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

func buildOpenApiMiddleware(
	logger *zap.SugaredLogger,
	schemaPath string,
) (func(next http.Handler) http.Handler, error) {
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
