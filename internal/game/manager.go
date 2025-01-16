package game

import (
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/mcoot/crosswordgame-go/internal/game/store"
	"github.com/mcoot/crosswordgame-go/internal/game/types"
)

type Manager struct {
	store store.GameStore
}

func NewGameManager(store store.GameStore) *Manager {
	return &Manager{
		store: store,
	}
}

func (m *Manager) NewGame(playerCount int) (types.GameId, error) {
	game := types.NewGame(playerCount)
	rawId, err := uuid.GenerateUUID()
	if err != nil {
		return "", err
	}
	id := types.GameId(rawId)
	err = m.store.StoreGame(id, game)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (m *Manager) GetGameState(gameId types.GameId) (*types.GameState, error) {
	game, err := m.store.RetrieveGame(gameId)
	if err != nil {
		return nil, err
	}
	return &game.GameState, nil
}

func (m *Manager) GetPlayerState(gameId types.GameId, playerId int) (*types.Player, error) {
	game, err := m.store.RetrieveGame(gameId)
	if err != nil {
		return nil, err
	}

	return getPlayer(game, playerId)
}

func (m *Manager) GetPlayerScore(gameId types.GameId, playerId int) (int, error) {
	game, err := m.store.RetrieveGame(gameId)
	if err != nil {
		return 0, err
	}

	player, err := getPlayer(game, playerId)
	if err != nil {
		return 0, err
	}

	return determineScore(player), nil
}

func getPlayer(game *types.Game, playerId int) (*types.Player, error) {
	if playerId < 0 || playerId >= len(game.Players) {
		return nil, fmt.Errorf("invalid player id %d", playerId)
	}

	return game.Players[playerId], nil
}
