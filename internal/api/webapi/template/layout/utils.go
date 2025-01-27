package layout

import (
	"context"
	"github.com/mcoot/crosswordgame-go/internal/api/webapi/template"
)

func shouldFullRerender(ctx context.Context) bool {
	renderCtx := template.GetRenderContext(ctx)
	if renderCtx != nil && renderCtx.Target.IsFullRefresh {
		return true
	}
	return false
}
