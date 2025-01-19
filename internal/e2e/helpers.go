package e2e

import (
	"github.com/mcoot/crosswordgame-go/internal/apitypes"
	"github.com/mcoot/crosswordgame-go/internal/client"
	"github.com/mcoot/crosswordgame-go/internal/game/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func createGame(t *testing.T, client *client.Client, playerCount int, boardDimension *int) types.GameId {
	t.Helper()

	createResp, err := client.CreateGame(playerCount, boardDimension)
	assert.NoError(t, err)
	assert.NotNil(t, createResp)
	assert.NotEmpty(t, createResp.GameId)

	return createResp.GameId
}

func getGameState(t *testing.T, client *client.Client, gameId types.GameId) *apitypes.GetGameStateResponse {
	t.Helper()

	gameState, err := client.GetGameState(gameId)
	assert.NoError(t, err)
	assert.NotNil(t, gameState)

	return gameState
}

func getPlayerState(t *testing.T, client *client.Client, gameId types.GameId, playerId int) *apitypes.GetPlayerStateResponse {
	t.Helper()

	playerState, err := client.GetPlayerState(gameId, playerId)
	assert.NoError(t, err)
	assert.NotNil(t, playerState)

	return playerState
}

func submitAnnouncement(t *testing.T, client *client.Client, gameId types.GameId, playerId int, letter string) {
	t.Helper()

	_, err := client.SubmitAnnouncement(gameId, playerId, letter)
	assert.NoError(t, err)
}

func submitPlacement(t *testing.T, client *client.Client, gameId types.GameId, playerId int, row, column int) {
	t.Helper()

	_, err := client.SubmitPlacement(gameId, playerId, row, column)
	assert.NoError(t, err)
}

func getPlayerScore(t *testing.T, client *client.Client, gameId types.GameId, playerId int) *apitypes.GetPlayerScoreResponse {
	t.Helper()

	playerScore, err := client.GetPlayerScore(gameId, playerId)
	assert.NoError(t, err)
	assert.NotNil(t, playerScore)

	return playerScore
}
