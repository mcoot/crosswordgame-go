package utils

import (
	"github.com/a-h/templ"
	"github.com/mcoot/crosswordgame-go/internal/api/webapi/rendering"
	"github.com/mcoot/crosswordgame-go/internal/api/webapi/template/common"
	"github.com/mcoot/crosswordgame-go/internal/api/webapi/template/pages"
	"github.com/mcoot/crosswordgame-go/internal/apitypes"
	"go.uber.org/zap"
	"net/http"
)

func PushUrl(w http.ResponseWriter, url string) {
	w.Header().Set("HX-Push-Url", url)
	w.Header().Add("Access-Control-Expose-Headers", "Hx-Push-Url")
}

func SendResponse(
	logger *zap.SugaredLogger,
	r *http.Request,
	w http.ResponseWriter,
	component templ.Component,
	code int,
) {
	htmx := rendering.GetHTMXProperties(r)
	renderCtx := rendering.WithRenderContext(r.Context(), &rendering.RenderContext{
		Target: rendering.RenderTarget{
			RefreshLevel:  htmx.DetermineRefreshLevel(),
			RefreshTarget: htmx.HTMXTarget,
		},
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)

	err := component.Render(renderCtx, w)
	if err != nil {
		logger.Errorw("error rendering web response", "error", err)
		return
	}
}

func Redirect(w http.ResponseWriter, r *http.Request, url string, code int) {
	http.Redirect(w, r, url, code)
}

func SendError(
	logger *zap.SugaredLogger,
	r *http.Request,
	w http.ResponseWriter,
	err error,
) {
	htmx := rendering.GetHTMXProperties(r)

	resp := apitypes.ToErrorResponse(err)
	logger.Warnw(
		"error handling web request",
		"message", resp.Message,
		"http_code", resp.HTTPCode,
		"kind", resp.Kind,
	)

	component := pages.Error(resp)
	// If the request is targeting a specific element that isn't the main page content, it should be displayed inline
	if htmx.DetermineRefreshLevel() == rendering.TargetedRefresh {
		component = common.ErrorInline(resp)
	}

	SendResponse(logger, r, w, component, resp.HTTPCode)
}
