package utils

import (
	"github.com/gorilla/mux"
	gametypes "github.com/mcoot/crosswordgame-go/internal/game/types"
	lobbytypes "github.com/mcoot/crosswordgame-go/internal/lobby/types"
	playertypes "github.com/mcoot/crosswordgame-go/internal/player/types"
	"net/http"
)

func GetGameIdPathParam(r *http.Request) gametypes.GameId {
	gameId, ok := mux.Vars(r)["gameId"]
	if !ok {
		return ""
	}
	return gametypes.GameId(gameId)
}

func GetPlayerIdPathParam(r *http.Request) playertypes.PlayerId {
	playerId, ok := mux.Vars(r)["playerId"]
	if !ok {
		return ""
	}
	return playertypes.PlayerId(playerId)
}

func GetLobbyIdPathParam(r *http.Request) lobbytypes.LobbyId {
	lobbyId, ok := mux.Vars(r)["lobbyId"]
	if !ok {
		return ""
	}
	return lobbytypes.LobbyId(lobbyId)
}
