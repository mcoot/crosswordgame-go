package client

import (
	"encoding/json"
	"github.com/mcoot/crosswordgame-go/internal/apitypes"
	"net/http"
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

func (c *Client) url(path string) string {
	return c.baseUrl + path
}
