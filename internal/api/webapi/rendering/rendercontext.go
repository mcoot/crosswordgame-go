package rendering

import (
	"context"
	"github.com/mcoot/crosswordgame-go/internal/utils"
)

type RenderContext struct {
	Target RenderTarget
}

func WithRenderContext(ctx context.Context, renderCtx *RenderContext) context.Context {
	return context.WithValue(ctx, utils.ContextKey("render_context"), renderCtx)
}

func GetRenderContext(ctx context.Context) *RenderContext {
	renderCtx, ok := ctx.Value(utils.ContextKey("render_context")).(*RenderContext)
	if !ok {
		return defaultRenderContext()
	}
	return renderCtx
}

func GetRenderRefreshLevel(ctx context.Context) RenderRefreshLevel {
	return GetRenderContext(ctx).Target.RefreshLevel
}

func defaultRenderContext() *RenderContext {
	return &RenderContext{
		Target: RenderTarget{
			RefreshLevel:  BrowserLevelRefresh,
			RefreshTarget: RefreshTargetNone,
		},
	}
}
