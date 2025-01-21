package webapi

import (
	"fmt"
	"github.com/a-h/templ"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	commonutils "github.com/mcoot/crosswordgame-go/internal/api/utils"
	"github.com/mcoot/crosswordgame-go/internal/api/webapi/template"
	"github.com/mcoot/crosswordgame-go/internal/api/webapi/utils"
	"github.com/mcoot/crosswordgame-go/internal/apitypes"
	"github.com/mcoot/crosswordgame-go/internal/errors"
	"github.com/mcoot/crosswordgame-go/internal/game"
	"github.com/mcoot/crosswordgame-go/internal/lobby"
	lobbytypes "github.com/mcoot/crosswordgame-go/internal/lobby/types"
	"github.com/mcoot/crosswordgame-go/internal/logging"
	"github.com/mcoot/crosswordgame-go/internal/player"
	playertypes "github.com/mcoot/crosswordgame-go/internal/player/types"
	"golang.org/x/tools/godoc/redirect"
	"net/http"
)

type CrosswordGameWebAPI struct {
	sessionStore  sessions.Store
	gameManager   *game.Manager
	lobbyManager  *lobby.Manager
	playerManager *player.Manager
}

func NewCrosswordGameWebAPI(
	sessionStore sessions.Store,
	gameManager *game.Manager,
	lobbyManager *lobby.Manager,
	playerManager *player.Manager,
) *CrosswordGameWebAPI {
	return &CrosswordGameWebAPI{
		sessionStore:  sessionStore,
		gameManager:   gameManager,
		lobbyManager:  lobbyManager,
		playerManager: playerManager,
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

	return nil
}

func (c *CrosswordGameWebAPI) Index(w http.ResponseWriter, r *http.Request) {
	session, err := commonutils.GetSessionDetails(c.sessionStore, r)
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

	var formComponent templ.Component
	if sessionIsLoggedIn {
		formComponent = templ.Join(template.LoggedInPlayerDetails(p), template.HostForm(), template.JoinForm())
	} else {
		formComponent = template.LoginForm()
	}

	utils.PushUrl(w, "/index")
	utils.SendResponse(logging.GetLogger(r.Context()), r, w, template.Index(formComponent), 200)
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

	session, err := commonutils.GetSessionDetails(c.sessionStore, r)
	if err != nil {
		utils.SendError(logging.GetLogger(r.Context()), r, w, err)
		return
	}

	session.PlayerId = playerId

	err = commonutils.SetSession(c.sessionStore, session, w, r)
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

	gameSpaceItem := template.EmptyGame()
	if lobbyState.RunningGame != nil {
		gameState, err := c.gameManager.GetGameState(lobbyState.RunningGame.GameId)
		if err != nil {
			utils.SendError(logging.GetLogger(r.Context()), r, w, err)
			return
		}

		gameSpaceItem = template.Game(gameState)
	}

	component := templ.Join(
		template.LobbyDetails(lobbyState),
		template.PlayerList(lobbyPlayers, player.Username),
		gameSpaceItem,
	)

	utils.PushUrl(w, fmt.Sprintf("/lobby/%s", lobbyId))
	utils.SendResponse(
		logging.GetLogger(r.Context()),
		r,
		w,
		component,
		200,
	)
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
		session, err := commonutils.GetSessionDetails(c.sessionStore, r)
		if err != nil {
			utils.SendError(logging.GetLogger(r.Context()), r, w, err)
			return
		}

		if session.PlayerId == "" {
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
