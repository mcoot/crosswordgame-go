package utils

import (
	"github.com/a-h/templ"
	"github.com/mcoot/crosswordgame-go/internal/api/webapi/template"
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
	renderCtx := template.WithRenderContext(r.Context(), &template.RenderContext{
		Target: template.RenderTarget{
			IsFullRefresh:  !htmx.IsTargeted(),
			SpecificTarget: "",
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
	if htmx.IsTargeted() {
		component = common.ErrorInline(resp)
	}

	SendResponse(logger, r, w, component, resp.HTTPCode)
}

type HTMXProperties struct {
	IsHTMX     bool
	IsBoosted  bool
	HTMXTarget string
}

func GetHTMXProperties(r *http.Request) HTMXProperties {
	return HTMXProperties{
		IsHTMX:     r.Header.Get("HX-Request") == "true",
		IsBoosted:  r.Header.Get("HX-Boosted") == "true",
		HTMXTarget: r.Header.Get("HX-Target"),
	}
}

func (p HTMXProperties) IsTargeted() bool {
	return p.IsHTMX && p.HTMXTarget != ""
}
