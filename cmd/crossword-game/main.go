package main

import (
	internalapi "github.com/mcoot/crosswordgame-go/internal/api"
	"github.com/mcoot/crosswordgame-go/internal/game"
	"github.com/mcoot/crosswordgame-go/internal/game/scoring"
	"github.com/mcoot/crosswordgame-go/internal/game/store"
	"github.com/mcoot/crosswordgame-go/internal/lobby"
	"github.com/mcoot/crosswordgame-go/internal/logging"
	"github.com/mcoot/crosswordgame-go/internal/utils"
	"log"
	"net/http"
)

func main() {
	ctx := utils.RootContext()
	ctx, err := logging.AddLoggerToContext(ctx, true)
	if err != nil {
		log.Fatalf("error adding logger to utils: %v", err)
	}
	logger := logging.GetLogger(ctx, "main")

	mux := http.NewServeMux()

	store := store.NewInMemoryStore()
	gameScorer, err := scoring.NewTxtDictScorer("./data/words.txt")
	if err != nil {
		logger.Fatalf("error loading dictionary: %v", err)
	}
	gameManager := game.NewGameManager(store, gameScorer)

	lobbyManager := lobby.NewLobbyManager(store)

	api := internalapi.NewCrosswordGameAPI(gameManager, lobbyManager)
	h, err := api.AttachToMux(ctx, mux, "./schema/openapi.yaml")
	if err != nil {
		logger.Fatalf("error setting up API: %v", err)
	}

	logger.Infow("starting server", "port", 8080)
	if err := http.ListenAndServe(":8080", h); err != nil {
		logger.Fatalf("error serving: %v", err)
	}
}
