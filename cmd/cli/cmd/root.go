package cmd

import (
	"context"
	"github.com/mcoot/crosswordgame-go/cmd/cli/cmd/game"
	"github.com/mcoot/crosswordgame-go/cmd/cli/cmd/lobby"
	"github.com/mcoot/crosswordgame-go/cmd/cli/cmd/player"
	"github.com/mcoot/crosswordgame-go/internal/cli"
	"github.com/mcoot/crosswordgame-go/internal/client"
	"github.com/mcoot/crosswordgame-go/internal/logging"
	"github.com/spf13/cobra"
	"net/http"
)

var (
	rootCmd = &cobra.Command{
		Use:   "crosswordgame",
		Short: "Crossword Game CLI",
		Long:  "CLI for interacting with the Crossword Game server",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			logger, err := logging.NewLogger(true)
			if err != nil {
				return err
			}
			ctx = logging.AddLoggerToContext(ctx, logger)

			cwgClient := initClient(cli.FlagServer)
			ctx = client.AddClientToContext(ctx, cwgClient)
			cmd.SetContext(ctx)

			return nil
		},
	}
)

func init() {
	cli.GlobalFlagServer(rootCmd)
	cli.GlobalFlagOutputMode(rootCmd)

	(&HealthCommand{}).Mount(rootCmd)
	(&game.GameCommand{}).Mount(rootCmd)
	(&lobby.LobbyCommand{}).Mount(rootCmd)
	(&player.PlayerCommand{}).Mount(rootCmd)
}

func initClient(baseUrl string) *client.Client {
	httpClient := &http.Client{}
	return client.NewClient(httpClient, baseUrl)
}

func Execute() error {
	return rootCmd.Execute()
}
