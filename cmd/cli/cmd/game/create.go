package game

import (
	"github.com/mcoot/crosswordgame-go/internal/cli"
	"github.com/mcoot/crosswordgame-go/internal/client"
	"github.com/spf13/cobra"
)

type CreateGameCommand struct {
	PlayerCount int
}

func (c *CreateGameCommand) Run(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	cwg := client.GetClient(ctx)

	game, err := cwg.CreateGame(c.PlayerCount)
	if err != nil {
		return err
	}

	return cli.WriteOutput(game)
}

func (c *CreateGameCommand) Mount(parent *cobra.Command) {
	createGameCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new game",
		Long:  "Create a new game",
		RunE:  c.Run,
	}

	createGameCmd.Flags().
		IntVarP(&c.PlayerCount, "players", "p", 2, "Number of players")
	_ = createGameCmd.MarkFlagRequired("players")

	parent.AddCommand(createGameCmd)
}
