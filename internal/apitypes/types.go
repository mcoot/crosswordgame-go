package apitypes

import (
	"fmt"
	gametypes "github.com/mcoot/crosswordgame-go/internal/game/types"
	lobbytypes "github.com/mcoot/crosswordgame-go/internal/lobby/types"
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
	GameId gametypes.GameId `json:"game_id"`
}

type GetGameStateResponse struct {
	Status                  gametypes.Status       `json:"status"`
	SquaresFilled           int                    `json:"squares_filled"`
	CurrentAnnouncingPlayer playertypes.PlayerId   `json:"current_announcing_player"`
	Players                 []playertypes.PlayerId `json:"players"`
}

type GetPlayerStateResponse struct {
	Board [][]string `json:"board"`
}

type GetPlayerScoreResponse struct {
	TotalScore int                     `json:"total_score"`
	Words      []*gametypes.ScoredWord `json:"words"`
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

type CreateLobbyRequest struct {
	Name string `json:"name"`
}

type CreateLobbyResponse struct {
	LobbyId lobbytypes.LobbyId `json:"lobby_id"`
}

type GetLobbyStateResponse struct {
	Name    string                 `json:"name"`
	Players []playertypes.PlayerId `json:"players"`
	GameID  gametypes.GameId       `json:"game_id,omitempty"`
}
