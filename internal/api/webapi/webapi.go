package webapi

import (
	"context"
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

func (c *CrosswordGameWebAPI) AttachToMux(ctx context.Context, mux *http.ServeMux) (http.Handler, error) {
	mux.Handle("GET /", redirect.Handler("/index.html"))
	mux.HandleFunc("GET /index.html", c.Index)

	return mux, nil
}

func (c *CrosswordGameWebAPI) Index(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello, World!"))
}
