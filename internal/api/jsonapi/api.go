package jsonapi

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mcoot/crosswordgame-go/internal/api/jsonapi/utils"
	commonutils "github.com/mcoot/crosswordgame-go/internal/api/utils"
	"github.com/mcoot/crosswordgame-go/internal/apitypes"
	"github.com/mcoot/crosswordgame-go/internal/game"
	"github.com/mcoot/crosswordgame-go/internal/lobby"
	"github.com/mcoot/crosswordgame-go/internal/logging"
	"go.uber.org/zap"
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

func (c *CrosswordGameAPI) AttachToRouter(router *mux.Router, baseLogger *zap.SugaredLogger, schemaPath string) error {
	err := setupMiddleware(router, baseLogger, schemaPath)
	if err != nil {
		return err
	}

	router.HandleFunc("/health", c.Healthcheck).Methods("GET")

	router.HandleFunc("/game", c.CreateGame).Methods("POST")
	router.HandleFunc("/game/{gameId}", c.GetGameState).Methods("GET")
	router.HandleFunc("/game/{gameId}/player/{playerId}", c.GetPlayerState).Methods("GET")
	router.HandleFunc("/game/{gameId}/player/{playerId}/announce", c.SubmitAnnouncement).Methods("POST")
	router.HandleFunc("/game/{gameId}/player/{playerId}/place", c.SubmitPlacement).Methods("POST")
	router.HandleFunc("/game/{gameId}/player/{playerId}/score", c.GetPlayerScore).Methods("GET")

	router.HandleFunc("/lobby", c.CreateLobby).Methods("POST")
	router.HandleFunc("/lobby/{lobbyId}", c.GetLobbyState).Methods("GET")
	router.HandleFunc("/lobby/{lobbyId}/join", c.JoinPlayerToLobby).Methods("POST")
	router.HandleFunc("/lobby/{lobbyId}/remove", c.RemovePlayerFromLobby).Methods("POST")
	router.HandleFunc("/lobby/{lobbyId}/attach", c.AttachGameToLobby).Methods("POST")
	router.HandleFunc("/lobby/{lobbyId}/detach", c.DetachGameFromLobby).Methods("POST")

	return nil
}

func (c *CrosswordGameAPI) Healthcheck(w http.ResponseWriter, r *http.Request) {
	utils.SendResponse(logging.GetLogger(r.Context()), w, &apitypes.HealthcheckResponse{
		Status:    "ok",
		StartTime: c.startTime.Format(time.RFC3339),
	}, 200)
}

func (c *CrosswordGameAPI) CreateGame(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger(r.Context())

	var req apitypes.CreateGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(logger, w, err)
		return
	}

	boardDimension := 5
	if req.BoardDimension != nil {
		boardDimension = *req.BoardDimension
	}

	gameId, err := c.gameManager.CreateGame(req.Players, boardDimension)
	if err != nil {
		utils.SendError(logger, w, err)
		return
	}

	w.Header().Add("Location", "/api/v1/game/"+string(gameId))

	utils.SendResponse(logger, w, apitypes.CreateGameResponse{GameId: gameId}, 201)
}

func (c *CrosswordGameAPI) GetGameState(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger(r.Context())
	gameId := commonutils.GetGameIdPathParam(r)

	gameState, err := c.gameManager.GetGameState(gameId)
	if err != nil {
		utils.SendError(logger, w, err)
		return
	}

	resp := apitypes.GetGameStateResponse{
		Status:                  gameState.Status,
		SquaresFilled:           gameState.SquaresFilled,
		CurrentAnnouncingPlayer: gameState.CurrentAnnouncingPlayer,
		Players:                 gameState.Players,
	}

	utils.SendResponse(logger, w, resp, 200)
}

func (c *CrosswordGameAPI) GetPlayerState(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger(r.Context())
	gameId := commonutils.GetGameIdPathParam(r)
	playerId := commonutils.GetPlayerIdPathParam(r)

	playerState, err := c.gameManager.GetPlayerBoard(gameId, playerId)
	if err != nil {
		utils.SendError(logger, w, err)
		return
	}

	resp := apitypes.GetPlayerStateResponse{
		Board: playerState.Data,
	}

	utils.SendResponse(logger, w, resp, 200)
}

func (c *CrosswordGameAPI) SubmitAnnouncement(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger(r.Context())
	gameId := commonutils.GetGameIdPathParam(r)
	playerId := commonutils.GetPlayerIdPathParam(r)

	var req apitypes.SubmitAnnouncementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(logger, w, err)
		return
	}

	err := c.gameManager.SubmitAnnouncement(gameId, playerId, req.Letter)
	if err != nil {
		utils.SendError(logger, w, err)
		return
	}

	utils.SendResponse(logger, w, &apitypes.SubmitAnnouncementResponse{}, 200)
}

func (c *CrosswordGameAPI) SubmitPlacement(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger(r.Context())
	gameId := commonutils.GetGameIdPathParam(r)
	playerId := commonutils.GetPlayerIdPathParam(r)

	var req apitypes.SubmitPlacementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(logger, w, err)
		return
	}

	err := c.gameManager.SubmitPlacement(gameId, playerId, req.Row, req.Column)
	if err != nil {
		utils.SendError(logger, w, err)
		return
	}

	utils.SendResponse(logger, w, &apitypes.SubmitPlacementResponse{}, 200)
}

func (c *CrosswordGameAPI) GetPlayerScore(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger(r.Context())
	gameId := commonutils.GetGameIdPathParam(r)
	playerId := commonutils.GetPlayerIdPathParam(r)

	totalScore, words, err := c.gameManager.GetPlayerScore(gameId, playerId)
	if err != nil {
		utils.SendError(logger, w, err)
		return
	}

	resp := apitypes.GetPlayerScoreResponse{
		TotalScore: totalScore,
		Words:      words,
	}

	utils.SendResponse(logger, w, resp, 200)
}

func (c *CrosswordGameAPI) CreateLobby(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger(r.Context())

	var req apitypes.CreateLobbyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(logger, w, err)
		return
	}

	lobbyId, err := c.lobbyManager.CreateLobby(req.Name)
	if err != nil {
		utils.SendError(logger, w, err)
		return
	}

	w.Header().Add("Location", "/api/v1/lobby/"+string(lobbyId))
	utils.SendResponse(logger, w, apitypes.CreateLobbyResponse{LobbyId: lobbyId}, 201)
}

func (c *CrosswordGameAPI) GetLobbyState(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger(r.Context())
	lobbyId := commonutils.GetLobbyIdPathParam(r)

	lobbyState, err := c.lobbyManager.GetLobbyState(lobbyId)
	if err != nil {
		utils.SendError(logger, w, err)
		return
	}

	resp := apitypes.GetLobbyStateResponse{
		Name:    lobbyState.Name,
		Players: lobbyState.Players,
	}
	if lobbyState.RunningGame != nil {
		resp.GameID = lobbyState.RunningGame.GameId
	}

	utils.SendResponse(logger, w, resp, 200)
}

func (c *CrosswordGameAPI) JoinPlayerToLobby(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger(r.Context())
	lobbyId := commonutils.GetLobbyIdPathParam(r)

	var req apitypes.JoinLobbyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(logger, w, err)
		return
	}

	err := c.lobbyManager.JoinPlayerToLobby(lobbyId, req.PlayerId)
	if err != nil {
		utils.SendError(logger, w, err)
		return
	}

	utils.SendResponse(logger, w, nil, 200)
}

func (c *CrosswordGameAPI) RemovePlayerFromLobby(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger(r.Context())
	lobbyId := commonutils.GetLobbyIdPathParam(r)

	var req apitypes.RemovePlayerFromLobbyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(logger, w, err)
		return
	}

	err := c.lobbyManager.RemovePlayerFromLobby(lobbyId, req.PlayerId)
	if err != nil {
		utils.SendError(logger, w, err)
		return
	}

	utils.SendResponse(logger, w, nil, 200)
}

func (c *CrosswordGameAPI) AttachGameToLobby(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger(r.Context())
	lobbyId := commonutils.GetLobbyIdPathParam(r)

	var req apitypes.AttachGameToLobbyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(logger, w, err)
		return
	}

	err := c.lobbyManager.AttachGameToLobby(lobbyId, req.GameId)
	if err != nil {
		utils.SendError(logger, w, err)
		return
	}

	utils.SendResponse(logger, w, nil, 200)
}

func (c *CrosswordGameAPI) DetachGameFromLobby(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger(r.Context())
	lobbyId := commonutils.GetLobbyIdPathParam(r)

	err := c.lobbyManager.DetachGameFromLobby(lobbyId)
	if err != nil {
		utils.SendError(logger, w, err)
		return
	}

	utils.SendResponse(logger, w, nil, 200)
}
