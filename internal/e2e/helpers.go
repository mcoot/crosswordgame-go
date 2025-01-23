package e2e

import (
	"github.com/mcoot/crosswordgame-go/internal/apitypes"
	"github.com/mcoot/crosswordgame-go/internal/client"
	gametypes "github.com/mcoot/crosswordgame-go/internal/game/types"
	lobbytypes "github.com/mcoot/crosswordgame-go/internal/lobby/types"
	playertypes "github.com/mcoot/crosswordgame-go/internal/player/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func createGame(t *testing.T, client *client.Client, players []playertypes.PlayerId, boardDimension *int) gametypes.GameId {
	t.Helper()

	createResp, err := client.CreateGame(players, boardDimension)
	assert.NoError(t, err)
	assert.NotNil(t, createResp)
	assert.NotEmpty(t, createResp.GameId)

	return createResp.GameId
}

func getGameState(t *testing.T, client *client.Client, gameId gametypes.GameId) *apitypes.GetGameStateResponse {
	t.Helper()

	gameState, err := client.GetGameState(gameId)
	assert.NoError(t, err)
	assert.NotNil(t, gameState)

	return gameState
}

func getPlayerState(t *testing.T, client *client.Client, gameId gametypes.GameId, playerId playertypes.PlayerId) *apitypes.GetPlayerStateResponse {
	t.Helper()

	playerState, err := client.GetPlayerState(gameId, playerId)
	assert.NoError(t, err)
	assert.NotNil(t, playerState)

	return playerState
}

func submitAnnouncement(t *testing.T, client *client.Client, gameId gametypes.GameId, playerId playertypes.PlayerId, letter string) {
	t.Helper()

	_, err := client.SubmitAnnouncement(gameId, playerId, letter)
	assert.NoError(t, err)
}

func submitPlacement(t *testing.T, client *client.Client, gameId gametypes.GameId, playerId playertypes.PlayerId, row, column int) {
	t.Helper()

	_, err := client.SubmitPlacement(gameId, playerId, row, column)
	assert.NoError(t, err)
}

func getPlayerScore(t *testing.T, client *client.Client, gameId gametypes.GameId, playerId playertypes.PlayerId) *apitypes.GetPlayerScoreResponse {
	t.Helper()

	playerScore, err := client.GetPlayerScore(gameId, playerId)
	assert.NoError(t, err)
	assert.NotNil(t, playerScore)

	return playerScore
}

func createLobby(t *testing.T, client *client.Client, name string) lobbytypes.LobbyId {
	t.Helper()

	createResp, err := client.CreateLobby(name)
	assert.NoError(t, err)
	assert.NotNil(t, createResp)
	assert.NotEmpty(t, createResp.LobbyId)

	return createResp.LobbyId
}

func getLobbyState(t *testing.T, client *client.Client, lobbyId lobbytypes.LobbyId) *apitypes.GetLobbyStateResponse {
	t.Helper()

	lobbyState, err := client.GetLobbyState(lobbyId)
	assert.NoError(t, err)
	assert.NotNil(t, lobbyState)

	return lobbyState
}

//func getLobbyForPlayer(t *testing.T, client *client.Client, playerId playertypes.PlayerId) *apitypes.GetLobbyStateResponse {
//	t.Helper()
//
//	lobby, err := client.GetLobbyForPlayer(playerId)
//	assert.NoError(t, err)
//	assert.NotNil(t, lobby)
//
//	return lobby
//}

func joinLobby(t *testing.T, client *client.Client, lobbyId lobbytypes.LobbyId, playerId playertypes.PlayerId) *apitypes.JoinLobbyResponse {
	t.Helper()

	resp, err := client.JoinLobby(lobbyId, playerId)
	assert.NoError(t, err)
	return resp
}

func removePlayerFromLobby(t *testing.T, client *client.Client, lobbyId lobbytypes.LobbyId, playerId playertypes.PlayerId) *apitypes.RemovePlayerFromLobbyResponse {
	t.Helper()

	resp, err := client.RemovePlayerFromLobby(lobbyId, playerId)
	assert.NoError(t, err)
	return resp
}

func attachGameToLobby(t *testing.T, client *client.Client, lobbyId lobbytypes.LobbyId, gameId gametypes.GameId) *apitypes.AttachGameToLobbyResponse {
	t.Helper()

	resp, err := client.AttachGameToLobby(lobbyId, gameId)
	assert.NoError(t, err)
	return resp
}

func detachGameFromLobby(t *testing.T, client *client.Client, lobbyId lobbytypes.LobbyId) *apitypes.DetachGameFromLobbyResponse {
	t.Helper()

	resp, err := client.DetachGameFromLobby(lobbyId)
	assert.NoError(t, err)
	return resp
}
