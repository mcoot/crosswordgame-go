package api

import (
	"context"
	"encoding/json"
	"github.com/mcoot/crosswordgame-go/internal/apitypes"
	"github.com/mcoot/crosswordgame-go/internal/errors"
	"github.com/mcoot/crosswordgame-go/internal/game"
	"github.com/mcoot/crosswordgame-go/internal/game/types"
	"github.com/mcoot/crosswordgame-go/internal/lobby"
	"github.com/mcoot/crosswordgame-go/internal/logging"
	"go.uber.org/zap"
	"net/http"
	"strconv"
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

	h, err := setupMiddleware(ctx, mux, schemaPath)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func (c *CrosswordGameAPI) Healthcheck(w http.ResponseWriter, r *http.Request) {
	logger := c.getLogger(r)

	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(apitypes.HealthcheckResponse{
		Status:    "ok",
		StartTime: c.startTime.Format(time.RFC3339),
	}); err != nil {
		logger.Errorw("error encoding response", "error", err)
		return
	}
}

func (c *CrosswordGameAPI) CreateGame(w http.ResponseWriter, r *http.Request) {
	logger := c.getLogger(r)

	var req apitypes.CreateGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.sendError(logger, w, err)
		return
	}

	boardDimension := 5
	if req.BoardDimension != nil {
		boardDimension = *req.BoardDimension
	}

	gameId, err := c.gameManager.NewGame(req.PlayerCount, boardDimension)
	if err != nil {
		c.sendError(logger, w, err)
		return
	}

	w.WriteHeader(201)
	if err := json.NewEncoder(w).Encode(apitypes.CreateGameResponse{GameId: gameId}); err != nil {
		logger.Errorw("error encoding response", "error", err)
		return
	}
}

func (c *CrosswordGameAPI) GetGameState(w http.ResponseWriter, r *http.Request) {
	logger := c.getLogger(r)
	gameId := getGameId(r)

	gameState, err := c.gameManager.GetGameState(gameId)
	if err != nil {
		c.sendError(logger, w, err)
		return
	}

	resp := apitypes.GetGameStateResponse{
		Status:                  gameState.Status,
		SquaresFilled:           gameState.SquaresFilled,
		CurrentAnnouncingPlayer: gameState.CurrentAnnouncingPlayer,
		PlayerCount:             gameState.PlayerCount,
	}

	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Errorw("error encoding response", "error", err)
		return
	}
}

func (c *CrosswordGameAPI) GetPlayerState(w http.ResponseWriter, r *http.Request) {
	logger := c.getLogger(r)
	gameId := getGameId(r)
	playerId, err := getPlayerId(r)
	if err != nil {
		c.sendError(logger, w, err)
		return
	}

	playerState, err := c.gameManager.GetPlayerState(gameId, playerId)
	if err != nil {
		c.sendError(logger, w, err)
		return
	}

	resp := apitypes.GetPlayerStateResponse{
		Board: playerState.Board.Data,
	}

	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Errorw("error encoding response", "error", err)
		return
	}
}

func (c *CrosswordGameAPI) SubmitAnnouncement(w http.ResponseWriter, r *http.Request) {
	logger := c.getLogger(r)
	gameId := getGameId(r)
	playerId, err := getPlayerId(r)
	if err != nil {
		c.sendError(logger, w, err)
		return
	}

	var req apitypes.SubmitAnnouncementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.sendError(logger, w, err)
		return
	}

	err = c.gameManager.SubmitAnnouncement(gameId, playerId, req.Letter)
	if err != nil {
		c.sendError(logger, w, err)
		return
	}

	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(&apitypes.SubmitAnnouncementResponse{}); err != nil {
		logger.Errorw("error encoding response", "error", err)
		return
	}
}

func (c *CrosswordGameAPI) SubmitPlacement(w http.ResponseWriter, r *http.Request) {
	logger := c.getLogger(r)
	gameId := getGameId(r)
	playerId, err := getPlayerId(r)
	if err != nil {
		c.sendError(logger, w, err)
		return
	}

	var req apitypes.SubmitPlacementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.sendError(logger, w, err)
		return
	}

	err = c.gameManager.SubmitPlacement(gameId, playerId, req.Row, req.Column)
	if err != nil {
		c.sendError(logger, w, err)
		return
	}

	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(&apitypes.SubmitPlacementResponse{}); err != nil {
		logger.Errorw("error encoding response", "error", err)
		return
	}
}

func (c *CrosswordGameAPI) GetPlayerScore(w http.ResponseWriter, r *http.Request) {
	logger := c.getLogger(r)
	gameId := getGameId(r)
	playerId, err := getPlayerId(r)
	if err != nil {
		c.sendError(logger, w, err)
		return
	}

	totalScore, words, err := c.gameManager.GetPlayerScore(gameId, playerId)
	if err != nil {
		c.sendError(logger, w, err)
		return
	}

	resp := apitypes.GetPlayerScoreResponse{
		TotalScore: totalScore,
		Words:      words,
	}

	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Errorw("error encoding response", "error", err)
		return
	}
}

func (c *CrosswordGameAPI) sendError(logger *zap.SugaredLogger, w http.ResponseWriter, err error) {
	var resp apitypes.ErrorResponse
	gameErr, ok := errors.AsGameError(err)
	if ok {
		resp = apitypes.ErrorResponse{
			Kind:     string(gameErr.Kind()),
			Message:  gameErr.Message(),
			HTTPCode: c.determineHttpErrorCode(gameErr),
		}
	} else {
		resp = apitypes.ErrorResponse{
			Kind:     "internal_error",
			Message:  err.Error(),
			HTTPCode: 500,
		}
	}

	logger.Warnw(
		"error handling request",
		"message", resp.Message,
		"http_code", resp.HTTPCode,
		"kind", resp.Kind,
	)
	w.WriteHeader(resp.HTTPCode)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Errorw("error encoding response", "error", err)
		return
	}
}

func (c *CrosswordGameAPI) determineHttpErrorCode(gameErr errors.GameError) int {
	switch gameErr.Kind() {
	case errors.GameErrorInvalidInput:
		return 400
	case errors.GameErrorNotFound:
		return 404
	case errors.GameErrorInvalidAction:
		return 400
	default:
		return 500
	}
}

func (c *CrosswordGameAPI) getLogger(r *http.Request) *zap.SugaredLogger {
	return logging.GetLogger(r.Context(), "api").
		With("path", r.URL.Path)
}

func getGameId(r *http.Request) types.GameId {
	return types.GameId(r.PathValue("gameId"))
}

func getPlayerId(r *http.Request) (int, error) {
	playerIdStr := r.PathValue("playerId")
	return strconv.Atoi(playerIdStr)
}
