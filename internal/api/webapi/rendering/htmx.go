package rendering

import "net/http"

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

func (p HTMXProperties) DetermineRefreshTarget() RenderRefreshTarget {
	if !p.IsHTMX {
		// The request isn't being made through ajax/htmx, so we need to send the whole document
		return RefreshTargetNone
	}
	return RenderRefreshTarget(p.HTMXTarget)
}
