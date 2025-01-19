package store

import (
	gametypes "github.com/mcoot/crosswordgame-go/internal/game/types"
	lobbytypes "github.com/mcoot/crosswordgame-go/internal/lobby/types"
)

type GameStore interface {
	StoreGame(gameId gametypes.GameId, game *gametypes.Game) error
	RetrieveGame(gameId gametypes.GameId) (*gametypes.Game, error)
}

type LobbyStore interface {
	StoreLobby(lobbyId lobbytypes.LobbyId, lobby *lobbytypes.Lobby) error
	RetrieveLobby(lobbyId lobbytypes.LobbyId) (*lobbytypes.Lobby, error)
}
