package main

import (
	"github.com/gorilla/sessions"
	"github.com/mcoot/crosswordgame-go/internal/api"
	"github.com/mcoot/crosswordgame-go/internal/logging"
	"github.com/mcoot/crosswordgame-go/internal/store"
	"log"
	"net/http"
	"os"

	_ "net/http/pprof"
)

func main() {
	if os.Getenv("PROFILE") == "true" {
		go func() {
			log.Print(http.ListenAndServe(":1234", nil))
		}()
	}

	logger, err := logging.NewLogger(true)
	if err != nil {
		log.Fatalf("error creating logger: %v", err)
	}

	logger.Infow("Initialising crossword-game")

	// TODO: Replace the key
	sessionStore := sessions.NewCookieStore([]byte("replace-me-key"))
	logger.Infow("Initialising datastore connection")
	db := store.NewInMemoryStore()
	handler, err := api.SetupAPI(
		logger,
		db,
		sessionStore,
		"./schema/openapi.yaml",
		"./data/words.txt",
	)
	if err != nil {
		logger.Fatalf("error setting up API: %v", err)
	}

	logger.Infow("starting server", "port", 8080)
	if err := http.ListenAndServe(":8080", handler); err != nil {
		logger.Fatalf("error serving: %v", err)
	}
}
