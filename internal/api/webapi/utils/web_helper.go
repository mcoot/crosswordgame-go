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
	htmx := GetHTMXProperties(r)
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
	htmx := GetHTMXProperties(r)

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

type HTMXProperties struct {
	IsHTMX     bool
	HTMXTarget string
}

func GetHTMXProperties(r *http.Request) HTMXProperties {
	return HTMXProperties{
		IsHTMX:     r.Header.Get("HX-Request") == "true",
		HTMXTarget: r.Header.Get("HX-Target"),
	}
}

func (p HTMXProperties) DetermineRefreshLevel() rendering.RenderRefreshLevel {
	if !p.IsHTMX {
		// The request isn't being made through ajax/htmx, so we need to send the whole document
		return rendering.BrowserLevelRefresh
	}
	if p.HTMXTarget == "" {
		// We aren't targeting anything, so we need to send the whole document
		return rendering.BrowserLevelRefresh
	}

	if rendering.RenderRefreshTarget(p.HTMXTarget) == rendering.RefreshTargetMain {
		// A page change will target the whole main div
		return rendering.PageChangeRefresh
	} else if rendering.RenderRefreshTarget(p.HTMXTarget) == rendering.RefreshTargetPageContent {
		// A change within one page will target the page content div
		return rendering.ContentRefresh
	} else {
		// Otherwise, some specific element is being targeted for change
		return rendering.TargetedRefresh
	}
}
