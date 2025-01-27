package layout

import (
	"context"
	"github.com/mcoot/crosswordgame-go/internal/api/webapi/rendering"
)

func getRenderLevel(ctx context.Context) rendering.RenderRefreshLevel {
	renderCtx := rendering.GetRenderContext(ctx)
	if renderCtx != nil {
		return renderCtx.Target.RefreshLevel
	}
	return rendering.BrowserLevelRefresh
}
