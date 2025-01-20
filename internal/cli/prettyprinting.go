package cli

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/mcoot/crosswordgame-go/internal/apitypes"
	"strings"
)

func PrettyPrint(v interface{}) {
	if ok := apiPrettyPrint(v); !ok {
		spew.Dump(v)
	}
}

func apiPrettyPrint(v interface{}) bool {
	switch v := v.(type) {
	case *apitypes.HealthcheckResponse:
		printHealthcheckResponse(v)
		return true
	case *apitypes.CreateGameResponse:
		printCreateGameResponse(v)
		return true
	case *apitypes.GetGameStateResponse:
		printGetGameStateResponse(v)
		return true
	case *apitypes.GetPlayerStateResponse:
		printGetPlayerStateResponse(v)
		return true
	case *apitypes.GetPlayerScoreResponse:
		printGetPlayerScoreResponse(v)
		return true
	case *apitypes.SubmitAnnouncementResponse:
		printSubmitAnnouncementResponse(v)
		return true
	case *apitypes.SubmitPlacementResponse:
		printSubmitPlacementResponse(v)
		return true
	case []interface{}:
		spew.Dump(v)
		return true
	}
	return false
}

func printHealthcheckResponse(v *apitypes.HealthcheckResponse) {
	fmt.Printf(`Server health:
  Start Time: %s
`, v.StartTime)
}

func printCreateGameResponse(v *apitypes.CreateGameResponse) {
	fmt.Printf(`Game created:
  Game ID: %s
`, v.GameId)
}

func printGetGameStateResponse(v *apitypes.GetGameStateResponse) {
	var playerSb strings.Builder
	for _, player := range v.Players {
		playerSb.WriteString(fmt.Sprintf("    %s\n", player))
	}

	fmt.Printf(`Game:
  Players:
%s
  Current State: %s
  Squares Filled: %d
  Current Announcing Player: %s
`, playerSb.String(), v.Status, v.SquaresFilled, v.CurrentAnnouncingPlayer)
}

func printGetPlayerStateResponse(v *apitypes.GetPlayerStateResponse) {
	fmt.Printf(`Player:
  Board:
`)
	printPlayerBoard(v.Board, 4)
}

func printSubmitAnnouncementResponse(v *apitypes.SubmitAnnouncementResponse) {
	fmt.Printf("Letter announced\n")
}

func printSubmitPlacementResponse(v *apitypes.SubmitPlacementResponse) {
	fmt.Printf("Letter placed\n")
}

func printPlayerBoard(board [][]string, indent int) {
	var sb strings.Builder

	for _, row := range board {
		sb.WriteString(strings.Repeat(" ", indent))
		for j, cell := range row {
			if cell == "" {
				sb.WriteString("∅")
			} else {
				sb.WriteString(cell)
			}
			if j < len(row)-1 {
				sb.WriteString(" ")
			}
		}
		sb.WriteString("\n")
	}

	fmt.Println(sb.String())
}

func printGetPlayerScoreResponse(v *apitypes.GetPlayerScoreResponse) {
	fmt.Printf(`Score:
  Total Score: %d
  Words:
`, v.TotalScore)
	for _, word := range v.Words {
		fmt.Printf("    %s (%d)\n", word.Word, word.Score)
	}
}
