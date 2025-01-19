package e2e

import (
	"context"
	internalapi "github.com/mcoot/crosswordgame-go/internal/api"
	"github.com/mcoot/crosswordgame-go/internal/client"
	"github.com/mcoot/crosswordgame-go/internal/game"
	"github.com/mcoot/crosswordgame-go/internal/game/scoring"
	"github.com/mcoot/crosswordgame-go/internal/game/store"
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
	gameStore := store.NewInMemoryStore()
	gameScorer, err := scoring.NewTxtDictScorer("../../data/words.txt")
	if err != nil {
		panic(err)
	}
	gameManager := game.NewGameManager(gameStore, gameScorer)
	api := internalapi.NewCrosswordGameAPI(gameManager)

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
