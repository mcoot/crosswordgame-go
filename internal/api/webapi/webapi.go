package webapi

import (
	"fmt"
	"github.com/a-h/templ"
	"github.com/gorilla/mux"
	commonutils "github.com/mcoot/crosswordgame-go/internal/api/utils"
	"github.com/mcoot/crosswordgame-go/internal/api/webapi/template"
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
	}
}

func (c *CrosswordGameWebAPI) AttachToRouter(router *mux.Router) error {
	router.NotFoundHandler = router.NewRoute().BuildOnly().Handler(NotFoundHandler()).GetHandler()

	router.Handle("/", redirect.Handler("/index")).Methods("GET")
	router.Handle("/index.html", redirect.Handler("/index")).Methods("GET")
	router.HandleFunc("/index", c.Index).Methods("GET")
	router.HandleFunc("/login", c.Login).Methods("POST")
	router.HandleFunc("/host", c.withLoggedInPlayer(c.StartLobbyAsHost)).Methods("POST")
	router.HandleFunc("/join", c.withLoggedInPlayer(c.JoinLobby)).Methods("POST")

	router.HandleFunc("/lobby/{lobbyId}", c.withLoggedInPlayer(c.LobbyPage)).Methods("GET")
	router.HandleFunc("/lobby/{lobbyId}/start", c.withLoggedInPlayer(c.StartNewGame)).Methods("POST")
	router.HandleFunc("/lobby/{lobbyId}/announce", c.withLoggedInPlayer(c.AnnounceLetter)).Methods("POST")
	router.HandleFunc("/lobby/{lobbyId}/place", c.withLoggedInPlayer(c.PlaceLetter)).Methods("POST")

	return nil
}

func (c *CrosswordGameWebAPI) Index(w http.ResponseWriter, r *http.Request) {
	session, err := c.sessionManager.GetSession(r)
	if err != nil {
		utils.SendError(logging.GetLogger(r.Context()), r, w, err)
		return
	}

	p, err := c.playerManager.LookupPlayer(session.PlayerId)
	sessionIsLoggedIn := true
	if err != nil {
		if errors.IsNotFoundError(err) {
			sessionIsLoggedIn = false
		} else {
			utils.SendError(logging.GetLogger(r.Context()), r, w, err)
			return
		}
	}

	var indexContents []templ.Component
	if sessionIsLoggedIn {
		indexContents = append(indexContents, template.LoggedInPlayerDetails(p))
		currentLobby, err := c.playerManager.GetLobbyForPlayer(p.Username)
		if err != nil {
			if errors.IsNotFoundError(err) {
				// Player not in a lobby
				indexContents = append(indexContents, template.NotInLobbyDetails())
			} else {
				utils.SendError(logging.GetLogger(r.Context()), r, w, err)
				return
			}
		} else {
			// Player in a lobby
			indexContents = append(indexContents, template.InLobbyDetails(currentLobby))
		}

	} else {
		indexContents = append(indexContents, template.LoginForm())
	}

	indexComponent := template.Index(templ.Join(indexContents...))
	utils.PushUrl(w, "/index")
	utils.SendResponse(logging.GetLogger(r.Context()), r, w, indexComponent, 200)
}

func (c *CrosswordGameWebAPI) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		utils.SendError(logging.GetLogger(r.Context()), r, w, err)
		return
	}

	displayName := r.PostForm.Get("display_name")
	if displayName == "" {
		utils.SendError(logging.GetLogger(r.Context()), r, w, fmt.Errorf("display_name is required"))
		return
	}

	playerId, err := c.playerManager.LoginAsEphemeral(displayName)
	if err != nil {
		utils.SendError(logging.GetLogger(r.Context()), r, w, err)
		return
	}

	session, err := c.sessionManager.GetSession(r)
	if err != nil {
		utils.SendError(logging.GetLogger(r.Context()), r, w, err)
		return
	}

	session.PlayerId = playerId

	err = c.sessionManager.SetSession(session, w, r)
	if err != nil {
		utils.SendError(logging.GetLogger(r.Context()), r, w, err)
		return
	}

	utils.Redirect(w, r, "/index", 303)
}

func (c *CrosswordGameWebAPI) StartLobbyAsHost(w http.ResponseWriter, r *http.Request, player *playertypes.Player) {
	err := r.ParseForm()
	if err != nil {
		utils.SendError(logging.GetLogger(r.Context()), r, w, err)
		return
	}

	lobbyName := r.PostForm.Get("lobby_name")
	if lobbyName == "" {
		utils.SendError(logging.GetLogger(r.Context()), r, w, fmt.Errorf("lobby_name is required"))
		return
	}

	lobbyId, err := c.lobbyManager.CreateLobby(lobbyName)
	if err != nil {
		utils.SendError(logging.GetLogger(r.Context()), r, w, err)
		return
	}

	err = c.lobbyManager.JoinPlayerToLobby(lobbyId, player.Username)
	if err != nil {
		// TODO: Scrap the lobby?
		utils.SendError(logging.GetLogger(r.Context()), r, w, err)
		return
	}

	utils.Redirect(w, r, fmt.Sprintf("/lobby/%s", lobbyId), 303)
}

func (c *CrosswordGameWebAPI) JoinLobby(w http.ResponseWriter, r *http.Request, player *playertypes.Player) {
	err := r.ParseForm()
	if err != nil {
		utils.SendError(logging.GetLogger(r.Context()), r, w, err)
		return
	}

	lobbyId := r.PostForm.Get("lobby_id")
	if lobbyId == "" {
		utils.SendError(logging.GetLogger(r.Context()), r, w, fmt.Errorf("lobby_id is required"))
		return
	}

	err = c.lobbyManager.JoinPlayerToLobby(lobbytypes.LobbyId(lobbyId), player.Username)
	if err != nil {
		utils.SendError(logging.GetLogger(r.Context()), r, w, err)
		return
	}

	utils.Redirect(w, r, fmt.Sprintf("/lobby/%s", lobbyId), 303)
}

func (c *CrosswordGameWebAPI) LobbyPage(w http.ResponseWriter, r *http.Request, player *playertypes.Player) {
	lobbyId := commonutils.GetLobbyIdPathParam(r)
	lobbyState, err := c.lobbyManager.GetLobbyState(lobbyId)
	if err != nil {
		utils.SendError(logging.GetLogger(r.Context()), r, w, err)
		return
	}

	lobbyPlayers := make([]*playertypes.Player, len(lobbyState.Players))
	for i, playerId := range lobbyState.Players {
		p, err := c.playerManager.LookupPlayer(playerId)
		if err != nil {
			utils.SendError(logging.GetLogger(r.Context()), r, w, err)
			return
		}

		lobbyPlayers[i] = p
	}

	gameComponents := []templ.Component{}

	if lobbyState.HasRunningGame() {
		gameState, err := c.gameManager.GetGameState(lobbyState.RunningGame.GameId)
		if err != nil {
			utils.SendError(logging.GetLogger(r.Context()), r, w, err)
			return
		}

		gamePlayers := make([]*playertypes.Player, len(gameState.Players))
		for i, playerId := range gameState.Players {
			p, err := c.playerManager.LookupPlayer(playerId)
			if err != nil {
				utils.SendError(logging.GetLogger(r.Context()), r, w, err)
				return
			}

			gamePlayers[i] = p
		}

		board, err := c.gameManager.GetPlayerBoard(gameState.Id, player.Username)
		if err != nil {
			utils.SendError(logging.GetLogger(r.Context()), r, w, err)
			return
		}

		canPlayerPlace := false
		if gameState.Status == gametypes.StatusAwaitingPlacement {
			hasPlayerPlaced, err := gameState.HasPlayerPlacedThisTurn(player.Username)
			if err != nil {
				utils.SendError(logging.GetLogger(r.Context()), r, w, err)
				return
			}
			if !hasPlayerPlaced {
				canPlayerPlace = true
			}
		}

		gameComponents = append(
			gameComponents,
			template.Game(
				gameState,
				gamePlayers,
				player,
				template.Board(lobbyId, player, board, canPlayerPlace),
			),
		)

		if gameState.Status == gametypes.StatusAwaitingAnnouncement &&
			gameState.CurrentAnnouncingPlayer == player.Username {
			gameComponents = append(gameComponents, template.AnnouncementForm(lobbyId))
		}

		if gameState.Status == gametypes.StatusFinished {
			gameComponents = append(
				gameComponents,
				template.GameScores(gamePlayers, player, gameState.PlayerScores),
				template.GameStartForm(lobbyId),
			)
		}
	} else {
		gameComponents = append(gameComponents, template.GameStartForm(lobbyId))
	}

	gameComponent := templ.Join(gameComponents...)
	component := template.Lobby(lobbyState, lobbyPlayers, player, gameComponent)
	utils.PushUrl(w, fmt.Sprintf("/lobby/%s", lobbyId))
	utils.SendResponse(
		logging.GetLogger(r.Context()),
		r,
		w,
		component,
		200,
	)
}

func (c *CrosswordGameWebAPI) StartNewGame(w http.ResponseWriter, r *http.Request, player *playertypes.Player) {
	lobbyId := commonutils.GetLobbyIdPathParam(r)
	lobbyState, err := c.lobbyManager.GetLobbyState(lobbyId)
	if err != nil {
		utils.SendError(logging.GetLogger(r.Context()), r, w, err)
		return
	}

	if lobbyState.HasRunningGame() {
		err = c.lobbyManager.DetachGameFromLobby(lobbyId)
		if err != nil {
			if errors.IsNotFoundError(err) {
				// The game was already detached, so we can just continue
			} else {
				utils.SendError(logging.GetLogger(r.Context()), r, w, err)
				return
			}
		}
	}

	err = r.ParseForm()
	if err != nil {
		utils.SendError(logging.GetLogger(r.Context()), r, w, err)
		return
	}

	boardSizeRaw := r.PostForm.Get("board_size")
	if boardSizeRaw == "" {
		utils.SendError(logging.GetLogger(r.Context()), r, w, fmt.Errorf("announced_letter is required"))
		return
	}

	boardSize, err := strconv.Atoi(boardSizeRaw)
	if err != nil {
		utils.SendError(logging.GetLogger(r.Context()), r, w, err)
		return
	}

	// TODO: Have some UX to handle players joining a lobby mid-game (who can't be in the game)
	gameId, err := c.gameManager.CreateGame(lobbyState.Players, boardSize)
	if err != nil {
		utils.SendError(logging.GetLogger(r.Context()), r, w, err)
		return
	}

	err = c.lobbyManager.AttachGameToLobby(lobbyId, gameId)
	if err != nil {
		utils.SendError(logging.GetLogger(r.Context()), r, w, err)
		return
	}

	utils.Redirect(w, r, fmt.Sprintf("/lobby/%s", lobbyId), 303)
}

func (c *CrosswordGameWebAPI) AnnounceLetter(w http.ResponseWriter, r *http.Request, player *playertypes.Player) {
	lobbyId := commonutils.GetLobbyIdPathParam(r)
	lobbyState, err := c.lobbyManager.GetLobbyState(lobbyId)
	if err != nil {
		utils.SendError(logging.GetLogger(r.Context()), r, w, err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		utils.SendError(logging.GetLogger(r.Context()), r, w, err)
		return
	}

	letter := r.PostForm.Get("announced_letter")
	if letter == "" {
		utils.SendError(logging.GetLogger(r.Context()), r, w, fmt.Errorf("announced_letter is required"))
		return
	}

	if !lobbyState.HasRunningGame() {
		utils.SendError(logging.GetLogger(r.Context()), r, w, &errors.InvalidActionError{
			Action: "place_letter",
			Reason: "the lobby has no running game",
		})
		return
	}

	err = c.gameManager.SubmitAnnouncement(lobbyState.RunningGame.GameId, player.Username, letter)
	if err != nil {
		utils.SendError(logging.GetLogger(r.Context()), r, w, err)
		return
	}

	utils.Redirect(w, r, fmt.Sprintf("/lobby/%s", lobbyId), 303)
}

func (c *CrosswordGameWebAPI) PlaceLetter(w http.ResponseWriter, r *http.Request, player *playertypes.Player) {
	lobbyId := commonutils.GetLobbyIdPathParam(r)
	lobbyState, err := c.lobbyManager.GetLobbyState(lobbyId)
	if err != nil {
		utils.SendError(logging.GetLogger(r.Context()), r, w, err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		utils.SendError(logging.GetLogger(r.Context()), r, w, err)
		return
	}

	rawRow := r.PostForm.Get("placement_row")
	if rawRow == "" {
		utils.SendError(logging.GetLogger(r.Context()), r, w, fmt.Errorf("placement_row is required"))
		return
	}
	row, err := strconv.Atoi(rawRow)
	if err != nil {
		utils.SendError(logging.GetLogger(r.Context()), r, w, err)
		return
	}

	rawColumn := r.PostForm.Get("placement_column")
	if rawColumn == "" {
		utils.SendError(logging.GetLogger(r.Context()), r, w, fmt.Errorf("placement_column is required"))
		return
	}
	column, err := strconv.Atoi(rawColumn)
	if err != nil {
		utils.SendError(logging.GetLogger(r.Context()), r, w, err)
		return
	}

	if !lobbyState.HasRunningGame() {
		utils.SendError(logging.GetLogger(r.Context()), r, w, &errors.InvalidActionError{
			Action: "place_letter",
			Reason: "the lobby has no running game",
		})
		return
	}

	err = c.gameManager.SubmitPlacement(lobbyState.RunningGame.GameId, player.Username, row, column)
	if err != nil {
		utils.SendError(logging.GetLogger(r.Context()), r, w, err)
		return
	}

	utils.Redirect(w, r, fmt.Sprintf("/lobby/%s", lobbyId), 303)
}

func NotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.SendError(logging.GetLogger(r.Context()), r, w, apitypes.ErrorResponse{
			HTTPCode: 404,
			Kind:     "not_found",
			Message:  "page not found",
		})
	})
}

func (c *CrosswordGameWebAPI) withLoggedInPlayer(
	f func(w http.ResponseWriter, r *http.Request, player *playertypes.Player),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := c.sessionManager.GetSession(r)
		if err != nil {
			utils.SendError(logging.GetLogger(r.Context()), r, w, err)
			return
		}

		if !session.IsLoggedIn() {
			utils.SendError(logging.GetLogger(r.Context()), r, w, apitypes.ErrorResponse{
				HTTPCode: 401,
				Kind:     "unauthorized",
				Message:  "not logged in",
			})
			return
		}

		p, err := c.playerManager.LookupPlayer(session.PlayerId)
		if err != nil {
			utils.SendError(logging.GetLogger(r.Context()), r, w, err)
			return
		}

		f(w, r, p)
	}
}
