package game

import (
	"fmt"
	"github.com/mcoot/crosswordgame-go/internal/errors"
	"github.com/mcoot/crosswordgame-go/internal/game/scoring"
	"github.com/mcoot/crosswordgame-go/internal/game/types"
	playertypes "github.com/mcoot/crosswordgame-go/internal/player/types"
	"github.com/mcoot/crosswordgame-go/internal/store"
	"slices"
	"strings"
)

type Manager struct {
	store  store.GameStore
	scorer scoring.Scorer
}

func NewGameManager(store store.GameStore, scorer scoring.Scorer) *Manager {
	return &Manager{
		store:  store,
		scorer: scorer,
	}
}

func (m *Manager) CreateGame(players []playertypes.PlayerId, boardDimension int) (types.GameId, error) {
	game, err := types.NewGame(players, boardDimension)
	if err != nil {
		return "", err
	}
	err = m.store.StoreGame(game.Id, game)
	if err != nil {
		return "", err
	}
	return game.Id, nil
}

func (m *Manager) GetGameState(gameId types.GameId) (*types.Game, error) {
	game, err := m.store.RetrieveGame(gameId)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (m *Manager) GetPlayerBoard(gameId types.GameId, playerId playertypes.PlayerId) (*types.Board, error) {
	game, err := m.store.RetrieveGame(gameId)
	if err != nil {
		return nil, err
	}

	return getPlayerBoard(game, playerId)
}

func (m *Manager) GetPlayerScore(gameId types.GameId, playerId playertypes.PlayerId) (int, []*types.ScoredWord, error) {
	game, err := m.store.RetrieveGame(gameId)
	if err != nil {
		return 0, nil, err
	}

	player, err := getPlayerBoard(game, playerId)
	if err != nil {
		return 0, nil, err
	}

	if game.Status != types.StatusFinished {
		return 0, nil, &errors.InvalidActionError{
			Action: "score",
			Reason: fmt.Sprintf(
				"game state is not %s, it is %s",
				types.StatusFinished,
				game.Status,
			),
		}
	}

	total, words := m.scorer.Score(player)
	return total, words, nil
}

func (m *Manager) SubmitAnnouncement(gameId types.GameId, playerId playertypes.PlayerId, announcedLetter string) error {
	game, err := m.store.RetrieveGame(gameId)
	if err != nil {
		return err
	}

	// Validate the player is real
	_, err = getPlayerBoard(game, playerId)
	if err != nil {
		return err
	}

	if game.Status != types.StatusAwaitingAnnouncement {
		return &errors.InvalidActionError{
			Action: "announce",
			Reason: fmt.Sprintf(
				"game state is not %s, it is %s",
				types.StatusAwaitingAnnouncement,
				game.Status,
			),
		}
	}
	if game.CurrentAnnouncingPlayer != playerId {
		return &errors.InvalidActionError{
			Action: "announce",
			Reason: fmt.Sprintf(
				"it is not player %s's turn to announce",
				playerId,
			),
		}
	}

	// Automatically upper-case the letter
	announcedLetter = strings.ToUpper(announcedLetter)

	if !types.IsValidLetter(announcedLetter) {
		return &errors.InvalidInputError{
			ErrMessage: fmt.Sprintf("invalid letter: %s", announcedLetter),
		}
	}

	game.Status = types.StatusAwaitingPlacement
	game.CurrentAnnouncedLetter = announcedLetter
	rotateAnnouncingPlayer(game)

	return m.store.StoreGame(gameId, game)
}

func (m *Manager) SubmitPlacement(gameId types.GameId, playerId playertypes.PlayerId, row, column int) error {
	game, err := m.store.RetrieveGame(gameId)
	if err != nil {
		return err
	}

	player, err := getPlayerBoard(game, playerId)
	if err != nil {
		return err
	}

	if game.Status != types.StatusAwaitingPlacement {
		return &errors.InvalidActionError{
			Action: "place",
			Reason: fmt.Sprintf(
				"game state is not %s, it is %s",
				types.StatusAwaitingPlacement,
				game.Status,
			),
		}
	}

	err = m.fillPlayerSquare(game, playerId, player, row, column)
	if err != nil {
		return err
	}

	err = m.checkAndProcessEndTurnOrGame(game)
	if err != nil {
		return err
	}

	return m.store.StoreGame(gameId, game)
}

func (m *Manager) fillPlayerSquare(
	game *types.Game,
	playerId playertypes.PlayerId,
	board *types.Board,
	row int,
	column int,
) error {
	playerFilledSquares := board.FilledSquares()
	if playerFilledSquares == game.SquaresFilled+1 {
		// The player already filled a square this turn
		return &errors.InvalidActionError{
			Action: "place",
			Reason: fmt.Sprintf("player %s has already placed a letter this turn", playerId),
		}
	}

	if playerFilledSquares != game.SquaresFilled {
		// Something went wrong in the game logic for them to not be on the correct # of squares
		// TODO: abandon game in this case
		return &errors.UnexpectedGameLogicError{
			ErrMessage: fmt.Sprintf(
				"expected player %s to have filled %d squares, but they have %d",
				playerId,
				game.SquaresFilled,
				playerFilledSquares,
			),
		}
	}

	if row < 0 || row >= board.Size() || column < 0 || column >= board.Size() {
		return &errors.InvalidInputError{
			ErrMessage: fmt.Sprintf("invalid row/column: %d/%d", row, column),
		}
	}

	if board.Data[row][column] != "" {
		return &errors.InvalidInputError{
			ErrMessage: fmt.Sprintf("square at row/column %d/%d is already filled", row, column),
		}
	}

	board.Data[row][column] = game.CurrentAnnouncedLetter
	return nil
}

func (m *Manager) checkAndProcessEndTurnOrGame(game *types.Game) error {
	// Check if any players are yet to have their turn
	playersLeft := false
	for _, board := range game.PlayerBoards {
		if board.FilledSquares() < game.SquaresFilled+1 {
			playersLeft = true
			break
		}
	}
	if playersLeft {
		return nil
	}

	// All players have had their turn, so end the round
	game.SquaresFilled++

	// Check if the game is over
	if game.SquaresFilled == game.TotalSquares() {
		game.Status = types.StatusFinished
	} else {
		// Proceed to the next turn
		game.Status = types.StatusAwaitingAnnouncement
	}
	return nil
}

func getPlayerBoard(game *types.Game, playerId playertypes.PlayerId) (*types.Board, error) {
	idx := slices.Index(game.Players, playerId)

	if idx == -1 {
		return nil, &errors.NotFoundError{
			ObjectKind: "player",
			ObjectID:   playerId,
		}
	}

	return game.PlayerBoards[idx], nil
}

func rotateAnnouncingPlayer(game *types.Game) {
	idx := slices.Index(game.Players, game.CurrentAnnouncingPlayer)
	game.CurrentAnnouncingPlayer = game.Players[(idx+1)%len(game.Players)]
}
