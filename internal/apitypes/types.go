package apitypes

import (
	"fmt"
	"github.com/mcoot/crosswordgame-go/internal/game/types"
)

type ErrorResponse struct {
	HTTPCode int    `json:"http_code"`
	Kind     string `json:"kind"`
	Message  string `json:"message"`
}

func (e ErrorResponse) Error() string {
	return fmt.Errorf("%d (%s): %s", e.HTTPCode, e.Kind, e.Message).Error()
}

type HealthcheckResponse struct {
	StartTime string `json:"start_time"`
}

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

type SubmitAnnouncementRequest struct {
	Letter string `json:"letter"`
}

type SubmitAnnouncementResponse struct{}

type SubmitPlacementRequest struct {
	Row    int `json:"row"`
	Column int `json:"column"`
}

type SubmitPlacementResponse struct{}
