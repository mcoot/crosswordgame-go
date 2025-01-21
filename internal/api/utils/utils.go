package utils

import (
	"github.com/gorilla/mux"
	"github.com/mcoot/crosswordgame-go/internal/game/types"
	types3 "github.com/mcoot/crosswordgame-go/internal/lobby/types"
	types2 "github.com/mcoot/crosswordgame-go/internal/player/types"
	"net/http"
)

func GetGameIdPathParam(r *http.Request) types.GameId {
	gameId, ok := mux.Vars(r)["gameId"]
	if !ok {
		return ""
	}
	return types.GameId(gameId)
}

func GetPlayerIdPathParam(r *http.Request) types2.PlayerId {
	playerId, ok := mux.Vars(r)["playerId"]
	if !ok {
		return ""
	}
	return types2.PlayerId(playerId)
}

func GetLobbyIdPathParam(r *http.Request) types3.LobbyId {
	lobbyId, ok := mux.Vars(r)["lobbyId"]
	if !ok {
		return ""
	}
	return types3.LobbyId(lobbyId)
}
