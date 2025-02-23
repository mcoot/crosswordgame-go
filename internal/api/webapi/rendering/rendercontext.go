package rendering

import (
	"context"
	lobbytypes "github.com/mcoot/crosswordgame-go/internal/lobby/types"
	playertypes "github.com/mcoot/crosswordgame-go/internal/player/types"
	"github.com/mcoot/crosswordgame-go/internal/utils"
)

type RenderContext struct {
	Target             RenderTarget
	LoggedInPlayer     *playertypes.Player
	CurrentPlayerLobby *lobbytypes.Lobby
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

func GetLoggedInPlayer(ctx context.Context) *playertypes.Player {
	return GetRenderContext(ctx).LoggedInPlayer
}

func GetCurrentPlayerLobby(ctx context.Context) *lobbytypes.Lobby {
	return GetRenderContext(ctx).CurrentPlayerLobby
}

func defaultRenderContext() *RenderContext {
	return &RenderContext{
		Target: RenderTarget{
			RefreshTarget: RefreshTargetNone,
		},
		LoggedInPlayer:     nil,
		CurrentPlayerLobby: nil,
	}
}
