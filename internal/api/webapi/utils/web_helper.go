package utils

import (
	"github.com/a-h/templ"
	"github.com/mcoot/crosswordgame-go/internal/api/webapi/template"
	"github.com/mcoot/crosswordgame-go/internal/apitypes"
	"go.uber.org/zap"
	"net/http"
)

func SendResponse(
	logger *zap.SugaredLogger,
	r *http.Request,
	w http.ResponseWriter,
	component templ.Component,
	code int,
) {
	htmx := GetHTMXRequestProperties(r)

	// If scripting is enabled and HTMX intends to swap out just the contents,
	// we don't need to re-send the layout, just the page contents
	// But for initial load and progressive enhancement,
	// we want to be able to swap out the whole page if necessary

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)

	var err error

	// If HTMX is present (HX-Request) and has a target (HX-Target), just render the page component
	if htmx.IsTargeted() {
		err = component.Render(r.Context(), w)
	} else {
		c := template.Layout(component)
		err = c.Render(r.Context(), w)
	}

	if err != nil {
		logger.Errorw("error rendering web response", "error", err)
		return
	}
}

func SendError(
	logger *zap.SugaredLogger,
	r *http.Request,
	w http.ResponseWriter,
	err error,
) {
	htmx := GetHTMXRequestProperties(r)

	resp := apitypes.ToErrorResponse(err)
	logger.Warnw(
		"error handling web request",
		"message", resp.Message,
		"http_code", resp.HTTPCode,
		"kind", resp.Kind,
	)

	component := template.ErrorPage(resp)
	if htmx.IsTargeted() {
		component = template.ErrorSpan(resp)
	}

	SendResponse(logger, r, w, component, resp.HTTPCode)
}

type HTMXRequestProperties struct {
	IsHTMX     bool
	HTMXTarget string
}

func GetHTMXRequestProperties(r *http.Request) HTMXRequestProperties {
	return HTMXRequestProperties{
		IsHTMX:     r.Header.Get("HX-Request") == "true",
		HTMXTarget: r.Header.Get("HX-Target"),
	}
}

func (p HTMXRequestProperties) IsTargeted() bool {
	return p.IsHTMX && p.HTMXTarget != ""
}
