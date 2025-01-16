package main

import (
	"github.com/mcoot/crosswordgame-go/cmd/cli/cmd"
	"github.com/mcoot/crosswordgame-go/internal/logging"
	"log"
)

func main() {
	logger, err := logging.NewLogger(true)
	if err != nil {
		log.Fatalf("Failed to create main logger: %v", err)
	}
	logger = logger.Named("main")

	if err := cmd.Execute(); err != nil {
		logger.Fatalw("Failed to execute command", "error", err)
	}
}
