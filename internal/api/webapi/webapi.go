package webapi

import (
	"fmt"
	"github.com/a-h/templ"
	"github.com/gorilla/mux"
	commonutils "github.com/mcoot/crosswordgame-go/internal/api/utils"
	"github.com/mcoot/crosswordgame-go/internal/api/webapi/rendering"
	gametemplates "github.com/mcoot/crosswordgame-go/internal/api/webapi/template/game"
	"github.com/mcoot/crosswordgame-go/internal/api/webapi/template/pages"
	"github.com/mcoot/crosswordgame-go/internal/api/webapi/utils"
	"github.com/mcoot/crosswordgame-go/internal/apitypes"
	"github.com/mcoot/crosswordgame-go/internal/errors"
	"github.com/mcoot/crosswordgame-go/internal/game"
	gametypes "github.com/mcoot/crosswordgame-go/internal/game/types"
	"github.com/mcoot/crosswordgame-go/internal/lobby"
	lobbytypes "github.com/mcoot/crosswordgame-go/internal/lobby/types"
	"github.com/mcoot/crosswordgame-go/internal/logging"
	"github.com/mcoot/crosswordgame-go/internal/player"
	playertypes "github.com/mcoot/crosswordgame-go/internal/player/types"
	"golang.org/x/tools/godoc/redirect"
	"net/http"
	"strconv"
)

type CrosswordGameWebAPI struct {
	sessionManager *commonutils.SessionManager
	gameManager    *game.Manager
	lobbyManager   *lobby.Manager
	playerManager  *player.Manager
	sseServer      *sseServer
}

func NewCrosswordGameWebAPI(
	sessionManager *commonutils.SessionManager,
	gameManager *game.Manager,
	lobbyManager *lobby.Manager,
	playerManager *player.Manager,
) *CrosswordGameWebAPI {
	return &CrosswordGameWebAPI{
		sessionManager: sessionManager,
		gameManager:    gameManager,
		lobbyManager:   lobbyManager,
		playerManager:  playerManager,
		sseServer:      newSSEServer(),
	}
}

func (c *CrosswordGameWebAPI) AttachToRouter(router *mux.Router) error {
	router.NotFoundHandler = router.NewRoute().BuildOnly().Handler(NotFoundHandler()).GetHandler()

	router.Use(c.sessionContextMiddleware)
	router.Use(renderContextMiddleware)

	router.Handle("/", redirect.Handler("/index")).Methods("GET")
	router.Handle("/index.html", redirect.Handler("/index")).Methods("GET")
	router.HandleFunc("/index", c.Index).Methods("GET")
	router.HandleFunc("/about", c.About).Methods("GET")
	router.HandleFunc("/login", c.Login).Methods("POST")
	router.HandleFunc("/logout", c.Logout).Methods("POST")
	router.HandleFunc("/host", c.StartLobbyAsHost).Methods("POST")
	router.HandleFunc("/join", c.JoinLobby).Methods("POST")

	router.HandleFunc("/lobby/{lobbyId}", c.LobbyPage).Methods("GET")
	router.HandleFunc("/lobby/{lobbyId}/leave", c.LeaveLobby).Methods("POST")
	router.HandleFunc("/lobby/{lobbyId}/start", c.StartNewGame).Methods("POST")
	router.HandleFunc("/lobby/{lobbyId}/abandon", c.AbandonGame).Methods("POST")
	router.HandleFunc("/lobby/{lobbyId}/announce", c.AnnounceLetter).Methods("POST")
	router.HandleFunc("/lobby/{lobbyId}/place", c.PlaceLetter).Methods("POST")

	c.sseServer.Start()
	router.HandleFunc("/lobby/{lobbyId}/sse/refresh", c.sseServer.HandleRequest).
		Methods("GET")

	return nil
}

func (c *CrosswordGameWebAPI) Index(w http.ResponseWriter, r *http.Request) {
	lobbyToJoin := lobbytypes.LobbyId(r.URL.Query().Get("join_lobby"))
	if lobbyToJoin != "" {
		_, err := c.lobbyManager.GetLobbyState(lobbyToJoin)
		if err != nil {
			utils.SendError(r, w, err)
			return
		}
	}

	indexComponent := pages.Index(lobbyToJoin)
	utils.PushUrl(w, "/index")
	utils.SendResponse(r, w, indexComponent, 200)
}

func (c *CrosswordGameWebAPI) About(w http.ResponseWriter, r *http.Request) {
	utils.PushUrl(w, "/about")
	utils.SendResponse(r, w, pages.About(), 200)
}

func (c *CrosswordGameWebAPI) Login(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger(r.Context())
	session, err := commonutils.GetSessionFromContext(r.Context())
	if err != nil {
		utils.SendError(r, w, err)
		return
	}
	if session.IsLoggedIn() {
		utils.SendError(r, w, &errors.InvalidActionError{
			Action: "login",
			Reason: "already logged in",
		})
		return
	}

	lobbyToJoin := lobbytypes.LobbyId(r.URL.Query().Get("join_lobby"))
	if lobbyToJoin != "" {
		_, err := c.lobbyManager.GetLobbyState(lobbyToJoin)
		if err != nil {
			utils.SendError(r, w, err)
			return
		}
	}

	err = r.ParseForm()
	if err != nil {
		utils.SendError(r, w, err)
		return
	}

	displayName := r.PostForm.Get("display_name")
	if displayName == "" {
		utils.SendError(r, w, fmt.Errorf("display_name is required"))
		return
	}

	playerId, err := c.playerManager.LoginAsEphemeral(displayName)
	if err != nil {
		utils.SendError(r, w, err)
		return
	}

	err = c.sessionManager.SaveLoggedInPlayer(w, r, playerId)
	if err != nil {
		utils.SendError(r, w, err)
		return
	}

	logger.Infow("player logged in", "player_id", playerId, "display_name", displayName)

	if lobbyToJoin != "" {
		err = c.lobbyManager.JoinPlayerToLobby(lobbyToJoin, playerId)
		if err != nil {
			utils.SendError(r, w, err)
			return
		}
		logger.Infow("player joined lobby", "lobby_id", lobbyToJoin, "player", playerId)
		c.sseServer.SendRefresh(lobbyToJoin, playerId)
		utils.Redirect(w, r, fmt.Sprintf("/lobby/%s", lobbyToJoin), 303)
		return
	}

	utils.Redirect(w, r, "/index", 303)
}

func (c *CrosswordGameWebAPI) Logout(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger(r.Context())
	session, err := getLoggedInSession(r)
	if err != nil {
		utils.SendError(r, w, err)
		return
	}

	if !session.IsLoggedIn() {
		utils.SendError(r, w, &errors.InvalidActionError{
			Action: "logout",
			Reason: "not currently logged in",
		})
		return
	}

	err = c.sessionManager.ClearLoggedInPlayer(w, r)
	if err != nil {
		utils.SendError(r, w, err)
		return
	}

	logger.Infow("player logged out", "player_id", session.Player.Username)

	utils.Redirect(w, r, "/index", 303)
}

func (c *CrosswordGameWebAPI) StartLobbyAsHost(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger(r.Context())
	session, err := getLoggedInSession(r)
	if err != nil {
		utils.SendError(r, w, err)
		return
	}

	if session.IsInLobby() {
		utils.SendError(r, w, &errors.InvalidActionError{
			Action: "start_lobby",
			Reason: "already in a lobby",
		})
		return
	}

	err = r.ParseForm()
	if err != nil {
		utils.SendError(r, w, err)
		return
	}

	lobbyName := r.PostForm.Get("lobby_name")
	if lobbyName == "" {
		utils.SendError(r, w, fmt.Errorf("lobby_name is required"))
		return
	}

	lobbyId, err := c.lobbyManager.CreateLobby(lobbyName)
	if err != nil {
		utils.SendError(r, w, err)
		return
	}

	err = c.lobbyManager.JoinPlayerToLobby(lobbyId, session.Player.Username)
	if err != nil {
		// TODO: Scrap the lobby?
		utils.SendError(r, w, err)
		return
	}

	logger.Infow(
		"lobby started",
		"lobby_id", lobbyId,
		"lobby_name", lobbyName,
		"hosting_player", session.Player.Username,
	)

	utils.Redirect(w, r, fmt.Sprintf("/lobby/%s", lobbyId), 303)
}

func (c *CrosswordGameWebAPI) JoinLobby(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger(r.Context())
	session, err := getLoggedInSession(r)
	if err != nil {
		utils.SendError(r, w, err)
		return
	}

	if session.IsInLobby() {
		utils.SendError(r, w, &errors.InvalidActionError{
			Action: "start_lobby",
			Reason: "already in a lobby",
		})
		return
	}

	err = r.ParseForm()
	if err != nil {
		utils.SendError(r, w, err)
		return
	}

	rawLobbyId := r.PostForm.Get("lobby_id")
	if rawLobbyId == "" {
		utils.SendError(r, w, fmt.Errorf("lobby_id is required"))
		return
	}

	lobbyId := lobbytypes.LobbyId(rawLobbyId)

	err = c.lobbyManager.JoinPlayerToLobby(lobbyId, session.Player.Username)
	if err != nil {
		utils.SendError(r, w, err)
		return
	}

	logger.Infow("player joined lobby", "lobby_id", lobbyId, "player", session.Player.Username)

	c.sseServer.SendRefresh(lobbyId, session.Player.Username)
	utils.Redirect(w, r, fmt.Sprintf("/lobby/%s", lobbyId), 303)
}

func (c *CrosswordGameWebAPI) LeaveLobby(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger(r.Context())
	session, err := getLoggedInSessionInLobby(r)
	if err != nil {
		utils.SendError(r, w, err)
		return
	}

	err = c.lobbyManager.RemovePlayerFromLobby(session.Lobby.Id, session.Player.Username)
	if err != nil {
		utils.SendError(r, w, err)
		return
	}

	logger.Infow("player left lobby", "lobby_id", session.Lobby.Id, "player", session.Player.Username)

	c.sseServer.SendRefresh(session.Lobby.Id, session.Player.Username)
	utils.Redirect(w, r, "/index", 303)
}

func (c *CrosswordGameWebAPI) LobbyPage(w http.ResponseWriter, r *http.Request) {
	session, err := getLoggedInSessionInLobby(r)
	if err != nil {
		utils.SendError(r, w, err)
		return
	}

	lobbyPlayers := make([]*playertypes.Player, len(session.Lobby.Players))
	for i, playerId := range session.Lobby.Players {
		p, err := c.playerManager.LookupPlayer(playerId)
		if err != nil {
			utils.SendError(r, w, err)
			return
		}

		lobbyPlayers[i] = p
	}

	var gameComponent templ.Component
	if session.Lobby.HasRunningGame() {
		var err error
		gameComponent, err = c.buildLobbyGameComponent(session.Player, session.Lobby)
		if err != nil {
			utils.SendError(r, w, err)
			return
		}
	} else {
		gameComponent = pages.GameStartForm(session.Lobby.Id)
	}

	component := pages.Lobby(session.Lobby, lobbyPlayers, session.Player, gameComponent)
	utils.PushUrl(w, fmt.Sprintf("/lobby/%s", session.Lobby.Id))
	utils.SendResponse(r, w, component, 200)
}

func (c *CrosswordGameWebAPI) buildLobbyGameComponent(
	player *playertypes.Player,
	lobbyState *lobbytypes.Lobby,
) (templ.Component, error) {
	gameState, err := c.gameManager.GetGameState(lobbyState.RunningGame.GameId)
	if err != nil {
		return nil, err
	}

	gamePlayers := make([]*playertypes.Player, len(gameState.Players))
	var currentAnnouncingPlayer *playertypes.Player
	for i, playerId := range gameState.Players {
		p, err := c.playerManager.LookupPlayer(playerId)
		if err != nil {
			return nil, err
		}

		gamePlayers[i] = p
		if playerId == gameState.CurrentAnnouncingPlayer {
			currentAnnouncingPlayer = p
		}
	}

	isPlayerInGame := false
	for _, playerId := range gameState.Players {
		if playerId == player.Username {
			isPlayerInGame = true
			break
		}
	}

	isGameFinished := gameState.Status == gametypes.StatusFinished

	gameStatusComponent := gametemplates.GameStatus(gameState, gamePlayers, currentAnnouncingPlayer, player, isPlayerInGame)

	var ingameComponent templ.Component
	if isPlayerInGame && !isGameFinished {
		ingameComponent, err = c.buildLobbyGamePlayerComponent(player, lobbyState, gameState, gamePlayers)
		if err != nil {
			return nil, err
		}
	} else {
		ingameComponent, err = c.buildLobbyGameSpectatorComponent(lobbyState, gameState, gamePlayers)
		if err != nil {
			return nil, err
		}
	}

	var components []templ.Component
	components = append(components, ingameComponent)
	if isGameFinished {
		components = append(
			components,
			gametemplates.GameScores(gamePlayers, player, gameState.PlayerScores),
			pages.GameStartForm(lobbyState.Id),
		)
	}
	components = append(components, pages.GameAbandonForm(lobbyState.Id, isGameFinished))

	return gametemplates.GameView(gameState, gamePlayers, player, gameStatusComponent, templ.Join(components...)), nil
}

func (c *CrosswordGameWebAPI) buildLobbyGamePlayerComponent(
	player *playertypes.Player,
	lobbyState *lobbytypes.Lobby,
	gameState *gametypes.Game,
	gamePlayers []*playertypes.Player,
) (templ.Component, error) {
	board, err := c.gameManager.GetPlayerBoard(gameState.Id, player.Username)
	if err != nil {
		return nil, err
	}

	canPlayerPlace := false
	if gameState.Status == gametypes.StatusAwaitingPlacement {
		hasPlayerPlaced, err := gameState.HasPlayerPlacedThisTurn(player.Username)
		if err != nil {
			return nil, err
		}
		if !hasPlayerPlaced {
			canPlayerPlace = true
		}
	}

	var components []templ.Component

	components = append(
		components,
		gametemplates.Board(lobbyState.Id, player, board, canPlayerPlace),
	)

	if gameState.Status == gametypes.StatusAwaitingAnnouncement &&
		gameState.CurrentAnnouncingPlayer == player.Username {
		components = append(components, gametemplates.AnnouncementForm(lobbyState.Id))
	}

	return templ.Join(components...), nil
}

func (c *CrosswordGameWebAPI) buildLobbyGameSpectatorComponent(
	lobbyState *lobbytypes.Lobby,
	gameState *gametypes.Game,
	gamePlayers []*playertypes.Player,
) (templ.Component, error) {
	var boardComponents []templ.Component
	for _, p := range gamePlayers {
		board, err := c.gameManager.GetPlayerBoard(gameState.Id, p.Username)
		if err != nil {
			return nil, err
		}
		boardComponents = append(boardComponents, gametemplates.Board(lobbyState.Id, p, board, false))
	}

	return templ.Join(boardComponents...), nil
}

func (c *CrosswordGameWebAPI) StartNewGame(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger(r.Context())
	session, err := getLoggedInSessionInLobby(r)
	if err != nil {
		utils.SendError(r, w, err)
		return
	}

	if session.Lobby.HasRunningGame() {
		err = c.lobbyManager.DetachGameFromLobby(session.Lobby.Id)
		if err != nil {
			if errors.IsNotFoundError(err) {
				// The game was already detached, so we can just continue
			} else {
				utils.SendError(r, w, err)
				return
			}
		}
	}

	err = r.ParseForm()
	if err != nil {
		utils.SendError(r, w, err)
		return
	}

	boardSizeRaw := r.PostForm.Get("board_size")
	if boardSizeRaw == "" {
		utils.SendError(r, w, fmt.Errorf("announced_letter is required"))
		return
	}

	boardSize, err := strconv.Atoi(boardSizeRaw)
	if err != nil {
		utils.SendError(r, w, err)
		return
	}

	if boardSize < 2 || boardSize > 10 {
		utils.SendError(r, w, fmt.Errorf("board_size must be between 2 and 10"))
		return
	}

	gameId, err := c.gameManager.CreateGame(session.Lobby.Players, boardSize)
	if err != nil {
		utils.SendError(r, w, err)
		return
	}

	err = c.lobbyManager.AttachGameToLobby(session.Lobby.Id, gameId)
	if err != nil {
		utils.SendError(r, w, err)
		return
	}

	logger.Infow(
		"game started",
		"game_id", gameId,
		"lobby_id", session.Lobby.Id,
		"board_size", boardSize,
	)

	c.sseServer.SendRefresh(session.Lobby.Id, session.Player.Username)
	utils.Redirect(w, r, fmt.Sprintf("/lobby/%s", session.Lobby.Id), 303)
}

func (c *CrosswordGameWebAPI) AbandonGame(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger(r.Context())
	session, err := getLoggedInSessionInLobby(r)
	if err != nil {
		utils.SendError(r, w, err)
		return
	}

	if !session.Lobby.HasRunningGame() {
		utils.SendError(r, w, &errors.InvalidActionError{
			Action: "abandon_game",
			Reason: "the lobby has no running game",
		})
		return
	}

	err = c.lobbyManager.DetachGameFromLobby(session.Lobby.Id)
	if err != nil {
		utils.SendError(r, w, err)
		return
	}

	logger.Infow("game cleared from lobby", "lobby_id", session.Lobby.Id)

	c.sseServer.SendRefresh(session.Lobby.Id, session.Player.Username)
	utils.Redirect(w, r, fmt.Sprintf("/lobby/%s", session.Lobby.Id), 303)
}

func (c *CrosswordGameWebAPI) AnnounceLetter(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger(r.Context())
	session, err := getLoggedInSessionInLobby(r)
	if err != nil {
		utils.SendError(r, w, err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		utils.SendError(r, w, err)
		return
	}

	letter := r.PostForm.Get("announced_letter")
	if letter == "" {
		utils.SendError(r, w, fmt.Errorf("announced_letter is required"))
		return
	}

	if !session.Lobby.HasRunningGame() {
		utils.SendError(r, w, &errors.InvalidActionError{
			Action: "place_letter",
			Reason: "the lobby has no running game",
		})
		return
	}

	err = c.gameManager.SubmitAnnouncement(session.Lobby.RunningGame.GameId, session.Player.Username, letter)
	if err != nil {
		utils.SendError(r, w, err)
		return
	}

	logger.Infow(
		"letter announced",
		"lobby_id", session.Lobby.Id,
		"player", session.Player.Username,
		"letter", letter,
	)

	c.sseServer.SendRefresh(session.Lobby.Id, session.Player.Username)
	utils.Redirect(w, r, fmt.Sprintf("/lobby/%s", session.Lobby.Id), 303)
}

func (c *CrosswordGameWebAPI) PlaceLetter(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger(r.Context())
	session, err := getLoggedInSessionInLobby(r)
	if err != nil {
		utils.SendError(r, w, err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		utils.SendError(r, w, err)
		return
	}

	rawRow := r.PostForm.Get("placement_row")
	if rawRow == "" {
		utils.SendError(r, w, fmt.Errorf("placement_row is required"))
		return
	}
	row, err := strconv.Atoi(rawRow)
	if err != nil {
		utils.SendError(r, w, err)
		return
	}

	rawColumn := r.PostForm.Get("placement_column")
	if rawColumn == "" {
		utils.SendError(r, w, fmt.Errorf("placement_column is required"))
		return
	}
	column, err := strconv.Atoi(rawColumn)
	if err != nil {
		utils.SendError(r, w, err)
		return
	}

	if !session.Lobby.HasRunningGame() {
		utils.SendError(r, w, &errors.InvalidActionError{
			Action: "place_letter",
			Reason: "the lobby has no running game",
		})
		return
	}

	err = c.gameManager.SubmitPlacement(session.Lobby.RunningGame.GameId, session.Player.Username, row, column)
	if err != nil {
		utils.SendError(r, w, err)
		return
	}

	logger.Infow(
		"letter placed",
		"lobby_id", session.Lobby.Id,
		"player", session.Player.Username,
		"row", row,
		"column", column,
	)

	c.sseServer.SendRefresh(session.Lobby.Id, session.Player.Username)
	utils.Redirect(w, r, fmt.Sprintf("/lobby/%s", session.Lobby.Id), 303)
}

func NotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.SendError(r, w, apitypes.ErrorResponse{
			HTTPCode: 404,
			Kind:     "not_found",
			Message:  "page not found",
		})
	})
}

func getLoggedInSession(r *http.Request) (*commonutils.Session, error) {
	session, err := commonutils.GetSessionFromContext(r.Context())
	if err != nil {
		return nil, err
	}

	if !session.IsLoggedIn() {
		return nil, apitypes.ErrorResponse{
			HTTPCode: 401,
			Kind:     "unauthorized",
			Message:  "not logged in",
		}
	}

	return session, nil
}

func getLoggedInSessionInLobby(r *http.Request) (*commonutils.Session, error) {
	session, err := getLoggedInSession(r)
	if err != nil {
		return nil, err
	}

	if !session.IsInLobby() {
		return nil, apitypes.ErrorResponse{
			HTTPCode: 400,
			Kind:     "bad_request",
			Message:  "not in a lobby",
		}
	}

	return session, nil
}

func (c *CrosswordGameWebAPI) sessionContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := c.sessionManager.GetSession(r, c.playerManager)
		if err != nil {
			utils.SendError(r, w, err)
			return
		}

		ctx := commonutils.AddSessionToContext(r.Context(), session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func renderContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := commonutils.GetSessionFromContext(r.Context())
		if err != nil {
			utils.SendError(r, w, err)
			return
		}

		htmx := rendering.GetHTMXProperties(r)
		renderCtx := &rendering.RenderContext{
			Target: rendering.RenderTarget{
				RefreshTarget: htmx.DetermineRefreshTarget(),
			},
			LoggedInPlayer:     session.Player,
			CurrentPlayerLobby: session.Lobby,
		}

		ctx := rendering.WithRenderContext(r.Context(), renderCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
