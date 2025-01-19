package store

import (
	gametypes "github.com/mcoot/crosswordgame-go/internal/game/types"
	lobbytypes "github.com/mcoot/crosswordgame-go/internal/lobby/types"
)

type InMemoryStore struct {
	games   map[gametypes.GameId]*gametypes.Game
	lobbies map[lobbytypes.LobbyId]*lobbytypes.Lobby
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		games: make(map[gametypes.GameId]*gametypes.Game),
	}
}

func (s *InMemoryStore) StoreGame(gameId gametypes.GameId, game *gametypes.Game) error {
	s.games[gameId] = game
	return nil
}

func (s *InMemoryStore) RetrieveGame(gameId gametypes.GameId) (*gametypes.Game, error) {
	game, ok := s.games[gameId]
	if !ok {
		return nil, &gametypes.NotFoundError{
			ObjectKind: "game",
			ObjectID:   gameId,
		}
	}
	return game, nil
}

func (s *InMemoryStore) StoreLobby(lobbyId lobbytypes.LobbyId, lobby *lobbytypes.Lobby) error {
	s.lobbies[lobbyId] = lobby
	return nil
}

func (s *InMemoryStore) RetrieveLobby(lobbyId lobbytypes.LobbyId) (*lobbytypes.Lobby, error) {
	lobby, ok := s.lobbies[lobbyId]
	if !ok {
		return nil, &gametypes.NotFoundError{
			ObjectKind: "lobby",
			ObjectID:   lobbyId,
		}
	}
	return lobby, nil
}
