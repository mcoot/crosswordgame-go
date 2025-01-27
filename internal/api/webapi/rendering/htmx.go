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

func (p HTMXProperties) DetermineRefreshLevel() RenderRefreshLevel {
	if !p.IsHTMX {
		// The request isn't being made through ajax/htmx, so we need to send the whole document
		return BrowserLevelRefresh
	}
	if p.HTMXTarget == "" {
		// We aren't targeting anything, so we need to send the whole document
		return BrowserLevelRefresh
	}

	if RenderRefreshTarget(p.HTMXTarget) == RefreshTargetMain {
		// A page change will target the whole main div
		return PageChangeRefresh
	} else if RenderRefreshTarget(p.HTMXTarget) == RefreshTargetPageContent {
		// A change within one page will target the page content div
		return ContentRefresh
	} else {
		// Otherwise, some specific element is being targeted for change
		return TargetedRefresh
	}
}
