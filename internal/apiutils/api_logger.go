package apiutils

import (
	"github.com/mcoot/crosswordgame-go/internal/logging"
	"go.uber.org/zap"
	"net/http"
)

func GetApiLogger(r *http.Request) *zap.SugaredLogger {
	return logging.GetLogger(r.Context(), "api").
		With("path", r.URL.Path)
}
