package main

import (
	"github.com/mcoot/crosswordgame-go/internal/api"
	"github.com/mcoot/crosswordgame-go/internal/logging"
	"github.com/mcoot/crosswordgame-go/internal/store"
	"log"
	"net/http"
)

func main() {
	logger, err := logging.NewLogger(true)
	if err != nil {
		log.Fatalf("error creating logger: %v", err)
	}
	db := store.NewInMemoryStore()
	router, err := api.SetupAPI(logger, db, "./schema/openapi.yaml", "./data/words.txt")
	if err != nil {
		logger.Fatalf("error setting up API: %v", err)
	}

	logger.Infow("starting server", "port", 8080)
	if err := http.ListenAndServe(":8080", router); err != nil {
		logger.Fatalf("error serving: %v", err)
	}
}
