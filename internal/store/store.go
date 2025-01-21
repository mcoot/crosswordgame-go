package store

import (
	gametypes "github.com/mcoot/crosswordgame-go/internal/game/types"
	lobbytypes "github.com/mcoot/crosswordgame-go/internal/lobby/types"
	playertypes "github.com/mcoot/crosswordgame-go/internal/player/types"
)

type GameStore interface {
	StoreGame(game *gametypes.Game) error
	RetrieveGame(gameId gametypes.GameId) (*gametypes.Game, error)
}

type LobbyStore interface {
	StoreLobby(lobby *lobbytypes.Lobby) error
	RetrieveLobby(lobbyId lobbytypes.LobbyId) (*lobbytypes.Lobby, error)
}

type PlayerStore interface {
	StorePlayer(player *playertypes.Player) error
	RetrievePlayer(playerId playertypes.PlayerId) (*playertypes.Player, error)
}

type Store interface {
	GameStore
	LobbyStore
	PlayerStore
}
