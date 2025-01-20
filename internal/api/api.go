package api

import (
	"context"
	"encoding/json"
	"github.com/mcoot/crosswordgame-go/internal/apitypes"
	"github.com/mcoot/crosswordgame-go/internal/apiutils"
	"github.com/mcoot/crosswordgame-go/internal/game"
	"github.com/mcoot/crosswordgame-go/internal/game/types"
	"github.com/mcoot/crosswordgame-go/internal/lobby"
	lobbytypes "github.com/mcoot/crosswordgame-go/internal/lobby/types"
	playertypes "github.com/mcoot/crosswordgame-go/internal/player/types"
	"net/http"
	"time"
)

type CrosswordGameAPI struct {
	startTime    time.Time
	gameManager  *game.Manager
	lobbyManager *lobby.Manager
}

func NewCrosswordGameAPI(gameManager *game.Manager, lobbyManager *lobby.Manager) *CrosswordGameAPI {
	return &CrosswordGameAPI{
		startTime:    time.Now(),
		gameManager:  gameManager,
		lobbyManager: lobbyManager,
	}
}

func (c *CrosswordGameAPI) AttachToMux(ctx context.Context, mux *http.ServeMux, schemaPath string) (http.Handler, error) {
	mux.Handle("GET /health", http.HandlerFunc(c.Healthcheck))

	mux.Handle("POST /api/v1/game", http.HandlerFunc(c.CreateGame))
	mux.Handle("GET /api/v1/game/{gameId}", http.HandlerFunc(c.GetGameState))
	mux.Handle("GET /api/v1/game/{gameId}/player/{playerId}", http.HandlerFunc(c.GetPlayerState))
	mux.Handle("POST /api/v1/game/{gameId}/player/{playerId}/announce", http.HandlerFunc(c.SubmitAnnouncement))
	mux.Handle("POST /api/v1/game/{gameId}/player/{playerId}/place", http.HandlerFunc(c.SubmitPlacement))
	mux.Handle("GET /api/v1/game/{gameId}/player/{playerId}/score", http.HandlerFunc(c.GetPlayerScore))

	mux.Handle("POST /api/v1/lobby", http.HandlerFunc(c.CreateLobby))
	mux.Handle("GET /api/v1/lobby/{lobbyId}", http.HandlerFunc(c.GetLobbyState))

	h, err := setupMiddleware(ctx, mux, schemaPath)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func (c *CrosswordGameAPI) Healthcheck(w http.ResponseWriter, r *http.Request) {
	apiutils.SendResponse(apiutils.GetApiLogger(r), w, &apitypes.HealthcheckResponse{
		Status:    "ok",
		StartTime: c.startTime.Format(time.RFC3339),
	}, 200)
}

func (c *CrosswordGameAPI) CreateGame(w http.ResponseWriter, r *http.Request) {
	logger := apiutils.GetApiLogger(r)

	var req apitypes.CreateGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apiutils.SendError(logger, w, err)
		return
	}

	boardDimension := 5
	if req.BoardDimension != nil {
		boardDimension = *req.BoardDimension
	}

	gameId, err := c.gameManager.NewGame(req.Players, boardDimension)
	if err != nil {
		apiutils.SendError(logger, w, err)
		return
	}

	w.Header().Add("Location", "/api/v1/game/"+string(gameId))

	apiutils.SendResponse(logger, w, apitypes.CreateGameResponse{GameId: gameId}, 201)
}

func (c *CrosswordGameAPI) GetGameState(w http.ResponseWriter, r *http.Request) {
	logger := apiutils.GetApiLogger(r)
	gameId := getGameId(r)

	gameState, err := c.gameManager.GetGameState(gameId)
	if err != nil {
		apiutils.SendError(logger, w, err)
		return
	}

	resp := apitypes.GetGameStateResponse{
		Status:                  gameState.Status,
		SquaresFilled:           gameState.SquaresFilled,
		CurrentAnnouncingPlayer: gameState.CurrentAnnouncingPlayer,
		Players:                 gameState.Players,
	}

	apiutils.SendResponse(logger, w, resp, 200)
}

func (c *CrosswordGameAPI) GetPlayerState(w http.ResponseWriter, r *http.Request) {
	logger := apiutils.GetApiLogger(r)
	gameId := getGameId(r)
	playerId := getPlayerId(r)

	playerState, err := c.gameManager.GetPlayerBoard(gameId, playerId)
	if err != nil {
		apiutils.SendError(logger, w, err)
		return
	}

	resp := apitypes.GetPlayerStateResponse{
		Board: playerState.Data,
	}

	apiutils.SendResponse(logger, w, resp, 200)
}

func (c *CrosswordGameAPI) SubmitAnnouncement(w http.ResponseWriter, r *http.Request) {
	logger := apiutils.GetApiLogger(r)
	gameId := getGameId(r)
	playerId := getPlayerId(r)

	var req apitypes.SubmitAnnouncementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apiutils.SendError(logger, w, err)
		return
	}

	err := c.gameManager.SubmitAnnouncement(gameId, playerId, req.Letter)
	if err != nil {
		apiutils.SendError(logger, w, err)
		return
	}

	apiutils.SendResponse(logger, w, &apitypes.SubmitAnnouncementResponse{}, 200)
}

func (c *CrosswordGameAPI) SubmitPlacement(w http.ResponseWriter, r *http.Request) {
	logger := apiutils.GetApiLogger(r)
	gameId := getGameId(r)
	playerId := getPlayerId(r)

	var req apitypes.SubmitPlacementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apiutils.SendError(logger, w, err)
		return
	}

	err := c.gameManager.SubmitPlacement(gameId, playerId, req.Row, req.Column)
	if err != nil {
		apiutils.SendError(logger, w, err)
		return
	}

	apiutils.SendResponse(logger, w, &apitypes.SubmitPlacementResponse{}, 200)
}

func (c *CrosswordGameAPI) GetPlayerScore(w http.ResponseWriter, r *http.Request) {
	logger := apiutils.GetApiLogger(r)
	gameId := getGameId(r)
	playerId := getPlayerId(r)

	totalScore, words, err := c.gameManager.GetPlayerScore(gameId, playerId)
	if err != nil {
		apiutils.SendError(logger, w, err)
		return
	}

	resp := apitypes.GetPlayerScoreResponse{
		TotalScore: totalScore,
		Words:      words,
	}

	apiutils.SendResponse(logger, w, resp, 200)
}

func (c *CrosswordGameAPI) CreateLobby(w http.ResponseWriter, r *http.Request) {
	logger := apiutils.GetApiLogger(r)

	var req apitypes.CreateLobbyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apiutils.SendError(logger, w, err)
		return
	}

	lobbyId, err := c.lobbyManager.NewLobby(req.Name)
	if err != nil {
		apiutils.SendError(logger, w, err)
		return
	}

	w.Header().Add("Location", "/api/v1/lobby/"+string(lobbyId))
	apiutils.SendResponse(logger, w, apitypes.CreateLobbyResponse{LobbyId: lobbyId}, 201)
}

func (c *CrosswordGameAPI) GetLobbyState(w http.ResponseWriter, r *http.Request) {
	logger := apiutils.GetApiLogger(r)
	lobbyId := getLobbyId(r)

	lobbyState, err := c.lobbyManager.GetLobbyState(lobbyId)
	if err != nil {
		apiutils.SendError(logger, w, err)
		return
	}

	resp := apitypes.GetLobbyStateResponse{
		Name:    lobbyState.Name,
		Players: lobbyState.Players,
	}
	if lobbyState.RunningGame != nil {
		resp.GameID = lobbyState.RunningGame.GameId
	}

	apiutils.SendResponse(logger, w, resp, 200)
}

func getGameId(r *http.Request) types.GameId {
	return types.GameId(r.PathValue("gameId"))
}

func getPlayerId(r *http.Request) playertypes.PlayerId {
	return playertypes.PlayerId(r.PathValue("playerId"))
}

func getLobbyId(r *http.Request) lobbytypes.LobbyId {
	return lobbytypes.LobbyId(r.PathValue("lobbyId"))
}
