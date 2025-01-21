package types

import (
	"github.com/hashicorp/go-uuid"
	"github.com/mcoot/crosswordgame-go/internal/errors"
	playertypes "github.com/mcoot/crosswordgame-go/internal/player/types"
	"slices"
)

type GameId string

type Status string

const (
	StatusAwaitingAnnouncement Status = "awaiting_announcement"
	StatusAwaitingPlacement    Status = "awaiting_placement"
	StatusFinished             Status = "finished"
)

type Game struct {
	Id                      GameId
	Status                  Status
	Players                 []playertypes.PlayerId
	SquaresFilled           int
	BoardDimension          int
	CurrentAnnouncingPlayer playertypes.PlayerId
	CurrentAnnouncedLetter  string
	PlayerBoards            map[playertypes.PlayerId]*Board
	PlayerScores            map[playertypes.PlayerId]*ScoreResult
}

func NewGame(players []playertypes.PlayerId, boardDimension int) (*Game, error) {
	rawId, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}
	id := GameId(rawId)

	playerBoards := make(map[playertypes.PlayerId]*Board)
	for _, p := range players {
		playerBoards[p] = NewBoard(boardDimension)
	}

	return &Game{
		Id:                      id,
		Status:                  StatusAwaitingAnnouncement,
		Players:                 players,
		SquaresFilled:           0,
		BoardDimension:          boardDimension,
		CurrentAnnouncingPlayer: players[0],
		CurrentAnnouncedLetter:  "",
		PlayerBoards:            playerBoards,
		PlayerScores:            make(map[playertypes.PlayerId]*ScoreResult),
	}, nil
}

func (g *Game) TotalSquares() int {
	return g.BoardDimension * g.BoardDimension
}

func (g *Game) GetIndexForPlayer(playerId playertypes.PlayerId) int {
	return slices.Index(g.Players, playerId)
}

func (g *Game) GetPlayerBoard(playerId playertypes.PlayerId) (*Board, error) {
	board, ok := g.PlayerBoards[playerId]
	if !ok {
		return nil, &errors.NotFoundError{
			ObjectKind: "player",
			ObjectID:   playerId,
		}
	}

	return board, nil
}

func (g *Game) GetPlayerScore(playerId playertypes.PlayerId) (*ScoreResult, error) {
	if g.GetIndexForPlayer(playerId) == -1 {
		return nil, &errors.NotFoundError{
			ObjectKind: "player",
			ObjectID:   playerId,
		}
	}

	score, ok := g.PlayerScores[playerId]
	if !ok {
		return nil, &errors.InvalidActionError{
			Action: "score",
			Reason: "player score is not yet calculated",
		}
	}

	return score, nil
}

func (g *Game) HasPlayerPlacedThisTurn(playerId playertypes.PlayerId) (bool, error) {
	board, err := g.GetPlayerBoard(playerId)
	if err != nil {
		return false, err
	}
	return g.SquaresFilled == board.FilledSquares()+1, nil
}
