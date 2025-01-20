package apiutils

import (
	"encoding/json"
	"github.com/mcoot/crosswordgame-go/internal/apitypes"
	"github.com/mcoot/crosswordgame-go/internal/errors"
	"go.uber.org/zap"
	"net/http"
)

func SendResponse(logger *zap.SugaredLogger, w http.ResponseWriter, resp interface{}, code int) {
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Errorw("error encoding response", "error", err)
		return
	}
}

func SendError(logger *zap.SugaredLogger, w http.ResponseWriter, err error) {
	var resp apitypes.ErrorResponse
	gameErr, ok := errors.AsGameError(err)
	if ok {
		resp = apitypes.ErrorResponse{
			Kind:     string(gameErr.Kind()),
			Message:  gameErr.Message(),
			HTTPCode: gameErr.HTTPCode(),
		}
	} else {
		resp = apitypes.ErrorResponse{
			Kind:     "internal_error",
			Message:  err.Error(),
			HTTPCode: 500,
		}
	}

	logger.Warnw(
		"error handling request",
		"message", resp.Message,
		"http_code", resp.HTTPCode,
		"kind", resp.Kind,
	)
	w.WriteHeader(resp.HTTPCode)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Errorw("error encoding response", "error", err)
		return
	}
}
