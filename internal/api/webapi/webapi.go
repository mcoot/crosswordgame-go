package webapi

import (
	"github.com/a-h/templ"
	"github.com/gorilla/mux"
	commonutils "github.com/mcoot/crosswordgame-go/internal/api/utils"
	"github.com/mcoot/crosswordgame-go/internal/api/webapi/template"
	"github.com/mcoot/crosswordgame-go/internal/api/webapi/utils"
	"github.com/mcoot/crosswordgame-go/internal/apitypes"
	"github.com/mcoot/crosswordgame-go/internal/game"
	"github.com/mcoot/crosswordgame-go/internal/lobby"
	"github.com/mcoot/crosswordgame-go/internal/logging"
	"golang.org/x/tools/godoc/redirect"
	"net/http"
)

type CrosswordGameWebAPI struct {
	gameManager  *game.Manager
	lobbyManager *lobby.Manager
}

func NewCrosswordGameWebAPI(gameManager *game.Manager, lobbyManager *lobby.Manager) *CrosswordGameWebAPI {
	return &CrosswordGameWebAPI{
		gameManager:  gameManager,
		lobbyManager: lobbyManager,
	}
}

func (c *CrosswordGameWebAPI) AttachToRouter(router *mux.Router) error {
	router.NotFoundHandler = router.NewRoute().BuildOnly().Handler(NotFoundHandler()).GetHandler()

	router.Handle("/", redirect.Handler("/index.html")).Methods("GET")
	router.Handle("/index", staticHandler(template.Index())).Methods("GET")
	router.HandleFunc("/lobby/{lobbyId}", c.LobbyPage).Methods("GET")

	return nil
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
		gameSpaceItem = template.Game(lobbyState.RunningGame.GameId)
	}

	utils.SendResponse(
		logging.GetLogger(r.Context()),
		r,
		w,
		template.LobbyPage(lobbyId, lobbyState, gameSpaceItem),
		200,
	)
}

func staticHandler(component templ.Component) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.SendResponse(logging.GetLogger(r.Context()), r, w, component, 200)
	})
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
