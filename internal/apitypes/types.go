package apitypes

import (
	"fmt"
	"github.com/mcoot/crosswordgame-go/internal/game/types"
	playertypes "github.com/mcoot/crosswordgame-go/internal/player/types"
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
	Status    string `json:"status"`
	StartTime string `json:"start_time"`
}

type CreateGameRequest struct {
	Players        []playertypes.PlayerId `json:"players"`
	BoardDimension *int                   `json:"board_dimension,omitempty"`
}

type CreateGameResponse struct {
	GameId types.GameId `json:"game_id"`
}

type GetGameStateResponse struct {
	Status                  types.Status           `json:"status"`
	SquaresFilled           int                    `json:"squares_filled"`
	CurrentAnnouncingPlayer playertypes.PlayerId   `json:"current_announcing_player"`
	Players                 []playertypes.PlayerId `json:"players"`
}

type GetPlayerStateResponse struct {
	Board [][]string `json:"board"`
}

type GetPlayerScoreResponse struct {
	TotalScore int                 `json:"total_score"`
	Words      []*types.ScoredWord `json:"words"`
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
