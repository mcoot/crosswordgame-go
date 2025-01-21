package webapi

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	commonutils "github.com/mcoot/crosswordgame-go/internal/api/utils"
	"github.com/mcoot/crosswordgame-go/internal/api/webapi/template"
	"github.com/mcoot/crosswordgame-go/internal/api/webapi/utils"
	"github.com/mcoot/crosswordgame-go/internal/apitypes"
	"github.com/mcoot/crosswordgame-go/internal/game"
	"github.com/mcoot/crosswordgame-go/internal/lobby"
	"github.com/mcoot/crosswordgame-go/internal/logging"
	"github.com/mcoot/crosswordgame-go/internal/player"
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
	router.HandleFunc("/index", c.Index).Methods("GET")
	router.HandleFunc("/lobby/{lobbyId}", c.LobbyPage).Methods("GET")

	return nil
}

func (c *CrosswordGameWebAPI) Index(w http.ResponseWriter, r *http.Request) {
	formComponent := template.LoginForm()

	utils.SendResponse(logging.GetLogger(r.Context()), r, w, template.Index(formComponent), 200)
}

func (c *CrosswordGameWebAPI) LobbyPage(w http.ResponseWriter, r *http.Request) {
	lobbyId := commonutils.GetLobbyIdPathParam(r)
	lobbyState, err := c.lobbyManager.GetLobbyState(lobbyId)
	if err != nil {
		utils.SendError(logging.GetLogger(r.Context()), r, w, err)
		return
	}

	gameSpaceItem := template.EmptyGameSpace()
	if lobbyState.RunningGame != nil {
		gameState, err := c.gameManager.GetGameState(lobbyState.RunningGame.GameId)
		if err != nil {
			utils.SendError(logging.GetLogger(r.Context()), r, w, err)
			return
		}

		gameSpaceItem = template.Game(gameState)
	}

	utils.SendResponse(
		logging.GetLogger(r.Context()),
		r,
		w,
		template.LobbyPage(lobbyState, gameSpaceItem),
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
