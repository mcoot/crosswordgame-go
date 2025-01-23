package store

import (
	"github.com/mcoot/crosswordgame-go/internal/errors"
	gametypes "github.com/mcoot/crosswordgame-go/internal/game/types"
	lobbytypes "github.com/mcoot/crosswordgame-go/internal/lobby/types"
	playertypes "github.com/mcoot/crosswordgame-go/internal/player/types"
)

type InMemoryStore struct {
	games   map[gametypes.GameId]*gametypes.Game
	lobbies map[lobbytypes.LobbyId]*lobbytypes.Lobby
	players map[playertypes.PlayerId]*playertypes.Player
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		games:   make(map[gametypes.GameId]*gametypes.Game),
		lobbies: make(map[lobbytypes.LobbyId]*lobbytypes.Lobby),
		players: make(map[playertypes.PlayerId]*playertypes.Player),
	}
}

func (s *InMemoryStore) StoreGame(game *gametypes.Game) error {
	s.games[game.Id] = game
	return nil
}

func (s *InMemoryStore) RetrieveGame(gameId gametypes.GameId) (*gametypes.Game, error) {
	game, ok := s.games[gameId]
	if !ok {
		return nil, &errors.NotFoundError{
			ObjectKind: "game",
			ObjectID:   gameId,
		}
	}
	return game, nil
}

func (s *InMemoryStore) StoreLobby(lobby *lobbytypes.Lobby) error {
	s.lobbies[lobby.Id] = lobby
	return nil
}

func (s *InMemoryStore) RetrieveLobby(lobbyId lobbytypes.LobbyId) (*lobbytypes.Lobby, error) {
	lobby, ok := s.lobbies[lobbyId]
	if !ok {
		return nil, &errors.NotFoundError{
			ObjectKind: "lobby",
			ObjectID:   lobbyId,
		}
	}
	return lobby, nil
}

func (s *InMemoryStore) StorePlayer(player *playertypes.Player) error {
	s.players[player.Username] = player
	return nil
}

func (s *InMemoryStore) RetrievePlayer(playerId playertypes.PlayerId) (*playertypes.Player, error) {
	player, ok := s.players[playerId]
	if !ok {
		return nil, &errors.NotFoundError{
			ObjectKind: "player",
			ObjectID:   playerId,
		}
	}
	return player, nil
}

func (s *InMemoryStore) RetrieveLobbyForPlayer(playerId playertypes.PlayerId) (*lobbytypes.Lobby, error) {
	// TODO: For an actual database, the DB layer should enforce the player being in one lobby
	for _, lobby := range s.lobbies {
		for _, playerIdInLobby := range lobby.Players {
			if playerIdInLobby == playerId {
				return lobby, nil
			}
		}
	}
	return nil, &errors.NotFoundError{
		ObjectKind: "lobby",
		KeyKind:    "player",
		ObjectID:   playerId,
	}
}
