package store

import (
	gametypes "github.com/mcoot/crosswordgame-go/internal/game/types"
	lobbytypes "github.com/mcoot/crosswordgame-go/internal/lobby/types"
	playertypes "github.com/mcoot/crosswordgame-go/internal/player/types"
)

type GameStore interface {
	StoreGame(gameId gametypes.GameId, game *gametypes.Game) error
	RetrieveGame(gameId gametypes.GameId) (*gametypes.Game, error)
}

type LobbyStore interface {
	StoreLobby(lobbyId lobbytypes.LobbyId, lobby *lobbytypes.Lobby) error
	RetrieveLobby(lobbyId lobbytypes.LobbyId) (*lobbytypes.Lobby, error)
}

type PlayerStore interface {
	StorePlayer(playerId playertypes.PlayerId, player *playertypes.Player) error
	RetrievePlayer(playerId playertypes.PlayerId) (*playertypes.Player, error)
}
