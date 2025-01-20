package webapi

import (
	"github.com/a-h/templ"
	"github.com/gorilla/mux"
	"github.com/mcoot/crosswordgame-go/internal/api/webapi/template"
	"github.com/mcoot/crosswordgame-go/internal/api/webapi/utils"
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
	router.Handle("/", redirect.Handler("/index.html")).Methods("GET")
	router.Handle("/index.html", staticHandler(template.Index())).Methods("GET")

	return nil
}

func staticHandler(component templ.Component) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.SendResponse(logging.GetLogger(r.Context()), r, w, component, 200)
	})
}
