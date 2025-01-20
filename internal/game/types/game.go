package types

import playertypes "github.com/mcoot/crosswordgame-go/internal/player/types"

type GameId string

type Status string

const (
	StatusAwaitingAnnouncement Status = "awaiting_announcement"
	StatusAwaitingPlacement    Status = "awaiting_placement"
	StatusFinished             Status = "finished"
)

type Game struct {
	Status                  Status
	Players                 []playertypes.PlayerId
	SquaresFilled           int
	BoardDimension          int
	CurrentAnnouncingPlayer playertypes.PlayerId
	CurrentAnnouncedLetter  string
	PlayerBoards            []*Board
}

func NewGame(players []playertypes.PlayerId, boardDimension int) *Game {
	playerBoards := make([]*Board, len(players))
	for i := range players {
		playerBoards[i] = NewBoard(boardDimension)
	}

	return &Game{
		Status:                  StatusAwaitingAnnouncement,
		Players:                 players,
		SquaresFilled:           0,
		BoardDimension:          boardDimension,
		CurrentAnnouncingPlayer: players[0],
		CurrentAnnouncedLetter:  "",
		PlayerBoards:            playerBoards,
	}
}

func (g *Game) TotalSquares() int {
	return g.BoardDimension * g.BoardDimension
}
