package rendering

import (
	"context"
	"github.com/mcoot/crosswordgame-go/internal/utils"
)

type RenderRefreshTarget string

const (
	RefreshTargetMain        RenderRefreshTarget = "#main"
	RefreshTargetPageContent RenderRefreshTarget = "#page-content"
)

type RenderRefreshLevel int

const (
	BrowserLevelRefresh RenderRefreshLevel = iota
	PageChangeRefresh
	ContentRefresh
	TargetedRefresh
)

type RenderTarget struct {
	RefreshLevel  RenderRefreshLevel
	RefreshTarget string
}

type RenderContext struct {
	Target RenderTarget
}

func WithRenderContext(ctx context.Context, renderCtx *RenderContext) context.Context {
	return context.WithValue(ctx, utils.ContextKey("render_context"), renderCtx)
}

func GetRenderContext(ctx context.Context) *RenderContext {
	renderCtx, _ := ctx.Value(utils.ContextKey("render_context")).(*RenderContext)
	return renderCtx
}

func GetRenderRefreshLevel(ctx context.Context) RenderRefreshLevel {
	renderCtx := GetRenderContext(ctx)
	if renderCtx != nil {
		return renderCtx.Target.RefreshLevel
	}
	return BrowserLevelRefresh
}
