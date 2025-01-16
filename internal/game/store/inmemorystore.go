package store

import (
	"github.com/mcoot/crosswordgame-go/internal/game/types"
)

type InMemoryStore struct {
	games map[types.GameId]*types.Game
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		games: make(map[types.GameId]*types.Game),
	}
}

func (s *InMemoryStore) StoreGame(gameId types.GameId, game *types.Game) error {
	s.games[gameId] = game
	return nil
}

func (s *InMemoryStore) RetrieveGame(gameId types.GameId) (*types.Game, error) {
	game, ok := s.games[gameId]
	if !ok {
		return nil, &types.NotFoundError{
			ObjectKind: "game",
			ObjectID:   gameId,
		}
	}
	return game, nil
}
