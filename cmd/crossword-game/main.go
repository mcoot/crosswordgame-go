package main

import (
	"github.com/gorilla/mux"
	internalapi "github.com/mcoot/crosswordgame-go/internal/api"
	"github.com/mcoot/crosswordgame-go/internal/game"
	"github.com/mcoot/crosswordgame-go/internal/game/scoring"
	"github.com/mcoot/crosswordgame-go/internal/lobby"
	"github.com/mcoot/crosswordgame-go/internal/logging"
	"github.com/mcoot/crosswordgame-go/internal/middleware"
	"github.com/mcoot/crosswordgame-go/internal/store"
	"log"
	"net/http"
)

func main() {
	logger, err := logging.NewLogger(true)
	if err != nil {
		log.Fatalf("error creating logger: %v", err)
	}

	router := mux.NewRouter()

	err = middleware.SetupMiddleware(router, logger)
	if err != nil {
		logger.Fatalf("error setting up middleware: %v", err)
	}

	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	db := store.NewInMemoryStore()
	gameScorer, err := scoring.NewTxtDictScorer("./data/words.txt")
	if err != nil {
		logger.Fatalf("error loading dictionary: %v", err)
	}
	gameManager := game.NewGameManager(db, gameScorer)

	lobbyManager := lobby.NewLobbyManager(db)

	api := internalapi.NewCrosswordGameAPI(gameManager, lobbyManager)
	err = api.AttachToRouter(apiRouter)
	if err != nil {
		logger.Fatalf("error setting up API: %v", err)
	}

	logger.Infow("starting server", "port", 8080)
	if err := http.ListenAndServe(":8080", router); err != nil {
		logger.Fatalf("error serving: %v", err)
	}
}
