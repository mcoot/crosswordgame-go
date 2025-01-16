package cmd

import (
	"context"
	"github.com/mcoot/crosswordgame-go/internal/cli"
	"github.com/mcoot/crosswordgame-go/internal/client"
	"github.com/mcoot/crosswordgame-go/internal/logging"
	"github.com/spf13/cobra"
	"net/http"
)

var (
	BaseUrl    string
	OutputMode cli.OutputMode

	rootCmd = &cobra.Command{
		Use:   "crosswordgame",
		Short: "Crossword Game CLI",
		Long:  "CLI for interacting with the Crossword Game server",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			ctx, err := logging.AddLoggerToContext(ctx, true)
			if err != nil {
				return err
			}

			cwgClient := initClient(BaseUrl)
			ctx = client.AddClientToContext(ctx, cwgClient)
			cmd.SetContext(ctx)

			return nil
		},
	}
)

func init() {
	rootCmd.PersistentFlags().
		StringVarP(&BaseUrl, "server", "s", "http://localhost:8080", "Server URL")
	rootCmd.PersistentFlags().
		VarP(&OutputMode, "output", "o", "Output mode (default: text, allowed: text, json, yaml)")

	rootCmd.AddCommand(healthCmd)
}

func initClient(baseUrl string) *client.Client {
	httpClient := &http.Client{}
	return client.NewClient(httpClient, baseUrl)
}

func Execute() error {
	return rootCmd.Execute()
}
