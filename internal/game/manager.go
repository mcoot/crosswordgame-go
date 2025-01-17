package game

import (
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/mcoot/crosswordgame-go/internal/game/store"
	"github.com/mcoot/crosswordgame-go/internal/game/types"
)

type Manager struct {
	store  store.GameStore
	scorer Scorer
}

func NewGameManager(store store.GameStore, scorer Scorer) *Manager {
	return &Manager{
		store:  store,
		scorer: scorer,
	}
}

func (m *Manager) NewGame(playerCount int, boardDimension int) (types.GameId, error) {
	game := types.NewGame(playerCount, boardDimension)
	rawId, err := uuid.GenerateUUID()
	if err != nil {
		return "", err
	}
	id := types.GameId(rawId)
	err = m.store.StoreGame(id, game)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (m *Manager) GetGameState(gameId types.GameId) (*types.GameState, error) {
	game, err := m.store.RetrieveGame(gameId)
	if err != nil {
		return nil, err
	}
	return &game.GameState, nil
}

func (m *Manager) GetPlayerState(gameId types.GameId, playerId int) (*types.Player, error) {
	game, err := m.store.RetrieveGame(gameId)
	if err != nil {
		return nil, err
	}

	return getPlayer(game, playerId)
}

func (m *Manager) GetPlayerScore(gameId types.GameId, playerId int) (int, []*types.ScoredWord, error) {
	game, err := m.store.RetrieveGame(gameId)
	if err != nil {
		return 0, nil, err
	}

	player, err := getPlayer(game, playerId)
	if err != nil {
		return 0, nil, err
	}

	if game.Status != types.StatusFinished {
		return 0, nil, &types.InvalidActionError{
			PlayerId: playerId,
			Action:   "score",
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

func (m *Manager) SubmitAnnouncement(gameId types.GameId, playerId int, announcedLetter string) error {
	game, err := m.store.RetrieveGame(gameId)
	if err != nil {
		return err
	}

	// Validate the player is real
	_, err = getPlayer(game, playerId)
	if err != nil {
		return err
	}

	if game.Status != types.StatusAwaitingAnnouncement {
		return &types.InvalidActionError{
			PlayerId: playerId,
			Action:   "announce",
			Reason: fmt.Sprintf(
				"game state is not %s, it is %s",
				types.StatusAwaitingAnnouncement,
				game.Status,
			),
		}
	}
	if game.CurrentAnnouncingPlayer != playerId {
		return &types.InvalidActionError{
			PlayerId: playerId,
			Action:   "announce",
			Reason: fmt.Sprintf(
				"it is not player %d's turn to announce",
				playerId,
			),
		}
	}
	if !types.IsValidLetter(announcedLetter) {
		return &types.InvalidInputError{
			ErrMessage: fmt.Sprintf("invalid letter: %s", announcedLetter),
		}
	}

	game.Status = types.StatusAwaitingPlacement
	game.CurrentAnnouncedLetter = announcedLetter
	game.CurrentAnnouncingPlayer = (playerId + 1) % len(game.Players)

	return m.store.StoreGame(gameId, game)
}

func (m *Manager) SubmitPlacement(gameId types.GameId, playerId int, row, column int) error {
	game, err := m.store.RetrieveGame(gameId)
	if err != nil {
		return err
	}

	player, err := getPlayer(game, playerId)
	if err != nil {
		return err
	}

	if game.Status != types.StatusAwaitingPlacement {
		return &types.InvalidActionError{
			PlayerId: playerId,
			Action:   "place",
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
	playerId int,
	player *types.Player,
	row int,
	column int,
) error {
	board := player.Board

	playerFilledSquares := board.FilledSquares()
	if playerFilledSquares == game.SquaresFilled+1 {
		// The player already filled a square this turn
		return &types.InvalidActionError{
			PlayerId: playerId,
			Action:   "place",
			Reason:   "player has already placed a letter this turn",
		}
	}

	if playerFilledSquares != game.SquaresFilled {
		// Something went wrong in the game logic for them to not be on the correct # of squares
		// TODO: abandon game in this case
		return &types.UnexpectedGameLogicError{
			ErrMessage: fmt.Sprintf(
				"expected player %d to have filled %d squares, but they have %d",
				playerId,
				game.SquaresFilled,
				playerFilledSquares,
			),
		}
	}

	if row < 0 || row >= board.Size() || column < 0 || column >= board.Size() {
		return &types.InvalidInputError{
			ErrMessage: fmt.Sprintf("invalid row/column: %d/%d", row, column),
		}
	}

	if board.Data[row][column] != "" {
		return &types.InvalidInputError{
			ErrMessage: fmt.Sprintf("square at row/column %d/%d is already filled", row, column),
		}
	}

	board.Data[row][column] = game.CurrentAnnouncedLetter
	return nil
}

func (m *Manager) checkAndProcessEndTurnOrGame(game *types.Game) error {
	// Check if any players are yet to have their turn
	playersLeft := false
	for _, player := range game.Players {
		if player.Board.FilledSquares() < game.SquaresFilled+1 {
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

func getPlayer(game *types.Game, playerId int) (*types.Player, error) {
	if playerId < 0 || playerId >= len(game.Players) {
		return nil, &types.NotFoundError{
			ObjectKind: "player",
			ObjectID:   playerId,
		}
	}

	return game.Players[playerId], nil
}
