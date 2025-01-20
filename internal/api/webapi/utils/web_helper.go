package utils

import (
	"github.com/a-h/templ"
	"github.com/mcoot/crosswordgame-go/internal/api/webapi/template"
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
	// If scripting is enabled and HTMX intends to swap out just the contents,
	// we don't need to re-send the layout, just the page contents
	// But for initial load and progressive enhancement,
	// we want to be able to swap out the whole page if necessary

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)

	var err error

	// If HTMX is present (HX-Request) and has a target (HX-Target), just render the page component
	if r.Header.Get("HX-Request") == "true" && r.Header.Get("HX-Target") != "" {
		err = component.Render(r.Context(), w)
	} else {
		c := template.Layout(component)
		err = c.Render(r.Context(), w)
	}

	if err != nil {
		logger.Errorw("error rendering response", "error", err)
		return
	}
}
