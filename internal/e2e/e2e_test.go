package e2e

import (
	"context"
	internalapi "github.com/mcoot/crosswordgame-go/internal/api"
	"github.com/mcoot/crosswordgame-go/internal/client"
	"github.com/mcoot/crosswordgame-go/internal/game"
	"github.com/mcoot/crosswordgame-go/internal/game/scoring"
	"github.com/mcoot/crosswordgame-go/internal/game/types"
	"github.com/mcoot/crosswordgame-go/internal/lobby"
	playertypes "github.com/mcoot/crosswordgame-go/internal/player/types"
	"github.com/mcoot/crosswordgame-go/internal/store"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type CrosswordGameE2ESuite struct {
	suite.Suite
	server *httptest.Server
	client *client.Client
}

func TestCrosswordGameE2ESuite(t *testing.T) {
	suite.Run(t, new(CrosswordGameE2ESuite))
}

func (s *CrosswordGameE2ESuite) SetupSuite() {
	store := store.NewInMemoryStore()
	gameScorer, err := scoring.NewTxtDictScorer("../../data/words.txt")
	if err != nil {
		panic(err)
	}
	gameManager := game.NewGameManager(store, gameScorer)

	lobbyManager := lobby.NewLobbyManager(store)

	api := internalapi.NewCrosswordGameAPI(gameManager, lobbyManager)

	mux := http.NewServeMux()
	h, err := api.AttachToMux(context.Background(), mux, "../../schema/openapi.yaml")
	if err != nil {
		panic(err)
	}

	// Run the API as an httptest server
	s.server = httptest.NewServer(h)
	s.client = client.NewClient(&http.Client{}, s.server.URL)
}

func (s *CrosswordGameE2ESuite) TearDownSuite() {
	s.server.Close()
}

func (s *CrosswordGameE2ESuite) Test_Healthcheck() {
	resp, err := s.client.Health()
	s.NoError(err)
	s.NotNil(resp)
	s.Equal("ok", resp.Status)
}

func (s *CrosswordGameE2ESuite) Test_FullGame2x2() {
	playerIds := []playertypes.PlayerId{
		"player0",
		"player1",
	}
	boardDim := 2
	// Create a game
	gameId := createGame(s.T(), s.client, playerIds, &boardDim)

	// Validate initial game state
	gameState := getGameState(s.T(), s.client, gameId)
	s.Equal(gameState.Players, playerIds)
	s.Equal(types.StatusAwaitingAnnouncement, gameState.Status)
	s.Equal(playerIds[0], gameState.CurrentAnnouncingPlayer)
	s.Equal(0, gameState.SquaresFilled)

	// Validate initial player state
	for _, playerId := range playerIds {
		playerState := getPlayerState(s.T(), s.client, gameId, playerId)
		s.Equal([][]string{{"", ""}, {"", ""}}, playerState.Board)
		// Getting score now should fail
		_, err := s.client.GetPlayerScore(gameId, playerId)
		s.Error(err)
	}

	// Player 1 attempting to announce should fail
	_, err := s.client.SubmitAnnouncement(gameId, playerIds[1], "a")
	s.Error(err)

	// Attempting to place now should fail
	_, err = s.client.SubmitPlacement(gameId, playerIds[0], 0, 0)
	s.Error(err)

	// Player 0 announces a letter
	submitAnnouncement(s.T(), s.client, gameId, playerIds[0], "a")

	// Validate game state after announcement
	gameState = getGameState(s.T(), s.client, gameId)
	s.Equal(types.StatusAwaitingPlacement, gameState.Status)
	s.Equal(playerIds[1], gameState.CurrentAnnouncingPlayer)
	s.Equal(0, gameState.SquaresFilled)

	// Both players place letters
	submitPlacement(s.T(), s.client, gameId, playerIds[0], 0, 0)
	submitPlacement(s.T(), s.client, gameId, playerIds[1], 1, 1)

	// Validate game state after placements
	gameState = getGameState(s.T(), s.client, gameId)
	s.Equal(types.StatusAwaitingAnnouncement, gameState.Status)
	s.Equal(playerIds[1], gameState.CurrentAnnouncingPlayer)
	s.Equal(1, gameState.SquaresFilled)

	// Validate player states after first round placements
	playerState0 := getPlayerState(s.T(), s.client, gameId, playerIds[0])
	s.Equal([][]string{{"A", ""}, {"", ""}}, playerState0.Board)
	playerState1 := getPlayerState(s.T(), s.client, gameId, playerIds[1])
	s.Equal([][]string{{"", ""}, {"", "A"}}, playerState1.Board)

	// Getting score now should still fail
	_, err = s.client.GetPlayerScore(gameId, playerIds[0])
	s.Error(err)
	_, err = s.client.GetPlayerScore(gameId, playerIds[1])
	s.Error(err)

	// Player 1 announces a letter
	submitAnnouncement(s.T(), s.client, gameId, playerIds[1], "s")

	// Validate game state after announcement
	gameState = getGameState(s.T(), s.client, gameId)
	s.Equal(types.StatusAwaitingPlacement, gameState.Status)
	s.Equal(playerIds[0], gameState.CurrentAnnouncingPlayer)
	s.Equal(1, gameState.SquaresFilled)

	// Player 0 attempting to overwrite an existing letter should fail
	_, err = s.client.SubmitPlacement(gameId, playerIds[0], 0, 0)
	s.Error(err)

	// Both players place letters
	submitPlacement(s.T(), s.client, gameId, playerIds[0], 1, 0)
	submitPlacement(s.T(), s.client, gameId, playerIds[1], 1, 0)

	// Validate game state after placements
	gameState = getGameState(s.T(), s.client, gameId)
	s.Equal(types.StatusAwaitingAnnouncement, gameState.Status)
	s.Equal(playerIds[0], gameState.CurrentAnnouncingPlayer)
	s.Equal(2, gameState.SquaresFilled)

	// Validate player states after second round placements
	playerState0 = getPlayerState(s.T(), s.client, gameId, playerIds[0])
	s.Equal([][]string{{"A", ""}, {"S", ""}}, playerState0.Board)
	playerState1 = getPlayerState(s.T(), s.client, gameId, playerIds[1])
	s.Equal([][]string{{"", ""}, {"S", "A"}}, playerState1.Board)

	// Play out the remaining two rounds
	submitAnnouncement(s.T(), s.client, gameId, playerIds[0], "t")
	submitPlacement(s.T(), s.client, gameId, playerIds[0], 0, 1)
	submitPlacement(s.T(), s.client, gameId, playerIds[1], 0, 1)
	submitAnnouncement(s.T(), s.client, gameId, playerIds[1], "e")
	submitPlacement(s.T(), s.client, gameId, playerIds[1], 0, 0)
	submitPlacement(s.T(), s.client, gameId, playerIds[0], 1, 1)

	// Validate game state after all placements
	gameState = getGameState(s.T(), s.client, gameId)
	s.Equal(types.StatusFinished, gameState.Status)
	s.Equal(4, gameState.SquaresFilled)

	// Validate the final game boards
	playerState0 = getPlayerState(s.T(), s.client, gameId, playerIds[0])
	s.Equal([][]string{{"A", "T"}, {"S", "E"}}, playerState0.Board)
	playerState1 = getPlayerState(s.T(), s.client, gameId, playerIds[1])
	s.Equal([][]string{{"E", "T"}, {"S", "A"}}, playerState1.Board)

	// Validate player scores
	playerScore0 := getPlayerScore(s.T(), s.client, gameId, playerIds[0])
	s.Equal([]*types.ScoredWord{
		{
			Word:        "AT",
			Score:       4,
			Direction:   types.ScoringDirectionHorizontal,
			StartRow:    0,
			StartColumn: 0,
		},
		{
			Word:        "AS",
			Score:       4,
			Direction:   types.ScoringDirectionVertical,
			StartRow:    0,
			StartColumn: 0,
		},
	}, playerScore0.Words)
	s.Equal(8, playerScore0.TotalScore)

	playerScore1 := getPlayerScore(s.T(), s.client, gameId, playerIds[1])
	s.Equal([]*types.ScoredWord{
		{
			Word:        "TA",
			Score:       4,
			Direction:   types.ScoringDirectionVertical,
			StartRow:    0,
			StartColumn: 1,
		},
	}, playerScore1.Words)
	s.Equal(4, playerScore1.TotalScore)

}
