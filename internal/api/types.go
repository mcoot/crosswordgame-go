package api

import "github.com/mcoot/crosswordgame-go/internal/game/types"

type CreateGameRequest struct {
	PlayerCount int `json:"player_count"`
}

type CreateGameResponse struct {
	GameId types.GameId `json:"game_id"`
}

type GetGameStateResponse struct {
	Status                  types.Status `json:"status"`
	SquaresFilled           int          `json:"squares_filled"`
	CurrentAnnouncingPlayer int          `json:"current_announcing_player"`
	PlayerCount             int          `json:"player_count"`
}

type GetPlayerStateResponse struct {
	Board [][]string `json:"board"`
}

type GetPlayerScoreResponse struct {
	Score int `json:"score"`
}
