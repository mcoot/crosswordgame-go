package utils

import (
	"encoding/json"
	"github.com/mcoot/crosswordgame-go/internal/apitypes"
	"go.uber.org/zap"
	"net/http"
)

func SendResponse(logger *zap.SugaredLogger, w http.ResponseWriter, resp interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Errorw("error encoding api response", "error", err)
		return
	}
}

func SendError(logger *zap.SugaredLogger, w http.ResponseWriter, err error) {
	resp := apitypes.ToErrorResponse(err)
	logger.Warnw(
		"error handling api request",
		"message", resp.Message,
		"http_code", resp.HTTPCode,
		"kind", resp.Kind,
	)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.HTTPCode)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Errorw("error encoding api response", "error", err)
		return
	}
}
