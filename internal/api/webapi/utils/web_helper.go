package utils

import (
	"github.com/a-h/templ"
	"github.com/mcoot/crosswordgame-go/internal/api/webapi/rendering"
	"github.com/mcoot/crosswordgame-go/internal/api/webapi/template/common"
	"github.com/mcoot/crosswordgame-go/internal/api/webapi/template/pages"
	"github.com/mcoot/crosswordgame-go/internal/apitypes"
	"github.com/mcoot/crosswordgame-go/internal/logging"
	"net/http"
)

func PushUrl(w http.ResponseWriter, url string) {
	w.Header().Set("HX-Push-Url", url)
	w.Header().Add("Access-Control-Expose-Headers", "Hx-Push-Url")
}

func SendResponse(
	r *http.Request,
	w http.ResponseWriter,
	component templ.Component,
	code int,
) {
	logger := logging.GetLogger(r.Context())

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)

	err := component.Render(r.Context(), w)
	if err != nil {
		logger.Errorw("error rendering web response", "error", err)
		return
	}
}

func Redirect(w http.ResponseWriter, r *http.Request, url string, code int) {
	http.Redirect(w, r, url, code)
}

func SendError(
	r *http.Request,
	w http.ResponseWriter,
	err error,
) {
	logger := logging.GetLogger(r.Context())

	resp := apitypes.ToErrorResponse(err)
	logger.Warnw(
		"error handling web request",
		"message", resp.Message,
		"http_code", resp.HTTPCode,
		"kind", resp.Kind,
	)

	component := pages.Error(resp)
	// If the request is targeting a specific element that isn't the main page content, it should be displayed inline
	if rendering.GetRenderContext(r.Context()).Target.RefreshLevel == rendering.TargetedRefresh {
		component = common.ErrorInline(resp)
	}

	SendResponse(r, w, component, resp.HTTPCode)
}
