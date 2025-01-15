package store

import "github.com/mcoot/crosswordgame-go/internal/game/types"

type GameStore interface {
	StoreGame(gameId types.GameId, game *types.Game) error
	RetrieveGame(gameId types.GameId) (*types.Game, error)
}
