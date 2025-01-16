package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mcoot/crosswordgame-go/internal/apitypes"
	"github.com/mcoot/crosswordgame-go/internal/game/types"
	"net/http"
)

const (
	createGamePath     = "/api/v1/game"
	getGameStatePath   = "/api/v1/game/%s"
	getPlayerStatePath = "/api/v1/game/%s/player/%d"
	getPlayerScorePath = "/api/v1/game/%s/player/%d/score"
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
	resp, err := c.client.Get(c.url("/health"))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var health apitypes.HealthcheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return nil, err
	}
	return &health, nil
}

func (c *Client) CreateGame(playerCount int) (*apitypes.CreateGameResponse, error) {
	body := apitypes.CreateGameRequest{
		PlayerCount: playerCount,
	}
	bodyJson, err := json.Marshal(body)

	resp, err := c.client.
		Post(c.url(createGamePath), "application/json", bytes.NewReader(bodyJson))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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

	var gameState apitypes.GetGameStateResponse
	if err := json.NewDecoder(resp.Body).Decode(&gameState); err != nil {
		return nil, err
	}
	return &gameState, nil
}

func (c *Client) GetPlayerState(gameId types.GameId, playerId int) (*apitypes.GetPlayerStateResponse, error) {
	resp, err := c.client.Get(c.url(fmt.Sprintf(getPlayerStatePath, gameId, playerId)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var playerState apitypes.GetPlayerStateResponse
	if err := json.NewDecoder(resp.Body).Decode(&playerState); err != nil {
		return nil, err
	}
	return &playerState, nil
}

func (c *Client) GetPlayerScore(gameId types.GameId, playerId int) (*apitypes.GetPlayerScoreResponse, error) {
	resp, err := c.client.Get(c.url(fmt.Sprintf(getPlayerScorePath, gameId, playerId)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var playerScore apitypes.GetPlayerScoreResponse
	if err := json.NewDecoder(resp.Body).Decode(&playerScore); err != nil {
		return nil, err
	}
	return &playerScore, nil
}

func (c *Client) url(path string) string {
	return c.baseUrl + path
}
