package api

import (
	"encoding/json"
	"github.com/mcoot/crosswordgame-go/internal/apitypes"
	"github.com/mcoot/crosswordgame-go/internal/game"
	"github.com/mcoot/crosswordgame-go/internal/game/types"
	"github.com/mcoot/crosswordgame-go/internal/logging"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

type CrosswordGameAPI struct {
	startTime   time.Time
	gameManager *game.Manager
}

func NewCrosswordGameAPI(gameManager *game.Manager) *CrosswordGameAPI {
	return &CrosswordGameAPI{
		startTime:   time.Now(),
		gameManager: gameManager,
	}
}

func (c *CrosswordGameAPI) AttachToMux(h *http.ServeMux) {
	h.Handle("GET /health", http.HandlerFunc(c.Healthcheck))
	h.Handle("POST /api/v1/game", http.HandlerFunc(c.CreateGame))
	h.Handle("GET /api/v1/game/{gameId}", http.HandlerFunc(c.GetGameState))
	h.Handle("GET /api/v1/game/{gameId}/player/{playerId}", http.HandlerFunc(c.GetPlayerState))
	h.Handle("POST /api/v1/game/{gameId}/player/{playerId}/announce", http.HandlerFunc(c.SubmitAnnouncement))
	h.Handle("POST /api/v1/game/{gameId}/player/{playerId}/place", http.HandlerFunc(c.SubmitPlacement))
	h.Handle("GET /api/v1/game/{gameId}/player/{playerId}/score", http.HandlerFunc(c.GetPlayerScore))
}

func (c *CrosswordGameAPI) Healthcheck(w http.ResponseWriter, r *http.Request) {
	logger := c.getLogger(r)

	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(apitypes.HealthcheckResponse{
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
		c.sendError(logger, w, 400, err)
		return
	}

	gameId, err := c.gameManager.NewGame(req.PlayerCount)
	if err != nil {
		c.sendError(logger, w, 500, err)
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
		// TODO: Appropriate error codes
		c.sendError(logger, w, 400, err)
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
		c.sendError(logger, w, 400, err)
		return
	}

	playerState, err := c.gameManager.GetPlayerState(gameId, playerId)
	if err != nil {
		// TODO: Appropriate error codes
		c.sendError(logger, w, 400, err)
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
		c.sendError(logger, w, 400, err)
		return
	}

	var req apitypes.SubmitAnnouncementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.sendError(logger, w, 400, err)
		return
	}

	err = c.gameManager.SubmitAnnouncement(gameId, playerId, req.Letter)
	if err != nil {
		// TODO: Appropriate error codes
		c.sendError(logger, w, 400, err)
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
		c.sendError(logger, w, 400, err)
		return
	}

	var req apitypes.SubmitPlacementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.sendError(logger, w, 400, err)
		return
	}

	err = c.gameManager.SubmitPlacement(gameId, playerId, req.Row, req.Column)
	if err != nil {
		// TODO: Appropriate error codes
		c.sendError(logger, w, 400, err)
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
		c.sendError(logger, w, 400, err)
		return
	}

	playerScore, err := c.gameManager.GetPlayerScore(gameId, playerId)
	if err != nil {
		// TODO: Appropriate error codes
		c.sendError(logger, w, 400, err)
		return
	}

	resp := apitypes.GetPlayerScoreResponse{
		Score: playerScore,
	}

	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Errorw("error encoding response", "error", err)
		return
	}
}

func (c *CrosswordGameAPI) sendError(logger *zap.SugaredLogger, w http.ResponseWriter, code int, err error) {
	logger.Warnw(
		"error handling request",
		"error", err,
		"code", code,
	)

	http.Error(w, err.Error(), code)
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
