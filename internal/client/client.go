package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mcoot/crosswordgame-go/internal/apitypes"
	"github.com/mcoot/crosswordgame-go/internal/game/types"
	lobbytypes "github.com/mcoot/crosswordgame-go/internal/lobby/types"
	playertypes "github.com/mcoot/crosswordgame-go/internal/player/types"
	"net/http"
)

const (
	healthcheckPath        = "/api/v1/health"
	createGamePath         = "/api/v1/game"
	getGameStatePath       = "/api/v1/game/%s"
	getPlayerStatePath     = "/api/v1/game/%s/player/%s"
	getPlayerScorePath     = "/api/v1/game/%s/player/%s/score"
	submitAnnouncementPath = "/api/v1/game/%s/player/%s/announce"
	submitPlacementPath    = "/api/v1/game/%s/player/%s/place"

	createLobbyPath         = "/api/v1/lobby"
	getLobbyStatePath       = "/api/v1/lobby/%s"
	joinLobbyPath           = "/api/v1/lobby/%s/join"
	removeFromLobbyPath     = "/api/v1/lobby/%s/remove"
	attachGameToLobbyPath   = "/api/v1/lobby/%s/attach"
	detachGameFromLobbyPath = "/api/v1/lobby/%s/detach"

	getLobbyForPlayerPath = "/api/v1/player/%s/lobby"
)

type Client struct {
	client  *http.Client
	baseUrl string
}

func NewClient(httpClient *http.Client, baseUrl string) *Client {
	return &Client{
		client:  httpClient,
		baseUrl: baseUrl,
	}
}

func (c *Client) Health() (*apitypes.HealthcheckResponse, error) {
	resp, err := c.client.Get(c.url(healthcheckPath))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, c.parseError(resp)
	}

	var health apitypes.HealthcheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return nil, err
	}
	return &health, nil
}

func (c *Client) CreateGame(players []playertypes.PlayerId, boardDimension *int) (*apitypes.CreateGameResponse, error) {
	body := apitypes.CreateGameRequest{
		Players: players,
	}
	if boardDimension != nil {
		body.BoardDimension = boardDimension
	}

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.
		Post(c.url(createGamePath), "application/json", bytes.NewReader(bodyJson))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return nil, c.parseError(resp)
	}

	var createGameResponse apitypes.CreateGameResponse
	if err := json.NewDecoder(resp.Body).Decode(&createGameResponse); err != nil {
		return nil, err
	}
	return &createGameResponse, nil
}

func (c *Client) GetGameState(gameId types.GameId) (*apitypes.GetGameStateResponse, error) {
	resp, err := c.client.Get(c.url(fmt.Sprintf(getGameStatePath, gameId)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, c.parseError(resp)
	}

	var gameState apitypes.GetGameStateResponse
	if err := json.NewDecoder(resp.Body).Decode(&gameState); err != nil {
		return nil, err
	}
	return &gameState, nil
}

func (c *Client) GetPlayerState(gameId types.GameId, playerId playertypes.PlayerId) (*apitypes.GetPlayerStateResponse, error) {
	resp, err := c.client.Get(c.url(fmt.Sprintf(getPlayerStatePath, gameId, playerId)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, c.parseError(resp)
	}

	var playerState apitypes.GetPlayerStateResponse
	if err := json.NewDecoder(resp.Body).Decode(&playerState); err != nil {
		return nil, err
	}
	return &playerState, nil
}

func (c *Client) GetPlayerScore(gameId types.GameId, playerId playertypes.PlayerId) (*apitypes.GetPlayerScoreResponse, error) {
	resp, err := c.client.Get(c.url(fmt.Sprintf(getPlayerScorePath, gameId, playerId)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, c.parseError(resp)
	}

	var playerScore apitypes.GetPlayerScoreResponse
	if err := json.NewDecoder(resp.Body).Decode(&playerScore); err != nil {
		return nil, err
	}
	return &playerScore, nil
}

func (c *Client) SubmitAnnouncement(
	gameId types.GameId,
	playerId playertypes.PlayerId,
	letter string,
) (*apitypes.SubmitAnnouncementResponse, error) {
	body := apitypes.SubmitAnnouncementRequest{
		Letter: letter,
	}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.
		Post(c.url(fmt.Sprintf(submitAnnouncementPath, gameId, playerId)), "application/json", bytes.NewReader(bodyJson))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, c.parseError(resp)
	}

	var ret apitypes.SubmitAnnouncementResponse
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (c *Client) SubmitPlacement(
	gameId types.GameId,
	playerId playertypes.PlayerId,
	row int,
	column int,
) (*apitypes.SubmitPlacementResponse, error) {
	body := apitypes.SubmitPlacementRequest{
		Row:    row,
		Column: column,
	}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.
		Post(c.url(fmt.Sprintf(submitPlacementPath, gameId, playerId)), "application/json", bytes.NewReader(bodyJson))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, c.parseError(resp)
	}

	var ret apitypes.SubmitPlacementResponse
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (c *Client) CreateLobby(name string) (*apitypes.CreateLobbyResponse, error) {
	body := apitypes.CreateLobbyRequest{
		Name: name,
	}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.
		Post(c.url(createLobbyPath), "application/json", bytes.NewReader(bodyJson))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return nil, c.parseError(resp)
	}

	var ret apitypes.CreateLobbyResponse
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (c *Client) GetLobbyState(lobbyId lobbytypes.LobbyId) (*apitypes.GetLobbyStateResponse, error) {
	resp, err := c.client.Get(c.url(fmt.Sprintf(getLobbyStatePath, lobbyId)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, c.parseError(resp)
	}

	var ret apitypes.GetLobbyStateResponse
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (c *Client) JoinLobby(lobbyId lobbytypes.LobbyId, playerId playertypes.PlayerId) (*apitypes.JoinLobbyResponse, error) {
	body := apitypes.JoinLobbyRequest{
		PlayerId: playerId,
	}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.
		Post(c.url(fmt.Sprintf(joinLobbyPath, lobbyId)), "application/json", bytes.NewReader(bodyJson))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, c.parseError(resp)
	}

	var ret apitypes.JoinLobbyResponse
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (c *Client) RemovePlayerFromLobby(lobbyId lobbytypes.LobbyId, playerId playertypes.PlayerId) (*apitypes.RemovePlayerFromLobbyResponse, error) {
	body := apitypes.RemovePlayerFromLobbyRequest{
		PlayerId: playerId,
	}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.
		Post(c.url(fmt.Sprintf(removeFromLobbyPath, lobbyId)), "application/json", bytes.NewReader(bodyJson))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, c.parseError(resp)
	}

	var ret apitypes.RemovePlayerFromLobbyResponse
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (c *Client) AttachGameToLobby(lobbyId lobbytypes.LobbyId, gameId types.GameId) (*apitypes.AttachGameToLobbyResponse, error) {
	body := apitypes.AttachGameToLobbyRequest{
		GameId: gameId,
	}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.
		Post(
			c.url(fmt.Sprintf(attachGameToLobbyPath, lobbyId)),
			"application/json",
			bytes.NewReader(bodyJson),
		)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, c.parseError(resp)
	}

	var ret apitypes.AttachGameToLobbyResponse
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (c *Client) DetachGameFromLobby(lobbyId lobbytypes.LobbyId) (*apitypes.DetachGameFromLobbyResponse, error) {
	body := apitypes.DetachGameFromLobbyRequest{}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.
		Post(
			c.url(fmt.Sprintf(detachGameFromLobbyPath, lobbyId)),
			"application/json",
			bytes.NewReader(bodyJson),
		)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, c.parseError(resp)
	}

	var ret apitypes.DetachGameFromLobbyResponse
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (c *Client) GetLobbyForPlayer(playerId playertypes.PlayerId) (*apitypes.GetLobbyStateResponse, error) {
	resp, err := c.client.Get(c.url(fmt.Sprintf(getLobbyForPlayerPath, playerId)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, c.parseError(resp)
	}

	var ret apitypes.GetLobbyStateResponse
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (c *Client) parseError(resp *http.Response) error {
	var apiErr apitypes.ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}
	return &apiErr
}

func (c *Client) url(path string) string {
	return c.baseUrl + path
}
