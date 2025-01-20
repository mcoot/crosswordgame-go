package webapi

import (
	"github.com/gorilla/mux"
	"github.com/mcoot/crosswordgame-go/internal/game"
	"github.com/mcoot/crosswordgame-go/internal/lobby"
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
	router.HandleFunc("/index.html", c.Index).Methods("GET")

	return nil
}

func (c *CrosswordGameWebAPI) Index(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello, World!"))
}
