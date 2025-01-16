package game

import (
	"github.com/mcoot/crosswordgame-go/internal/cli"
	"github.com/spf13/cobra"
)

type GetGameStateCommand struct {
	GameId string
}

func (c *GetGameStateCommand) Run(cmd *cobra.Command, args []string) error {
	return nil
}

func (c *GetGameStateCommand) Mount(parent *cobra.Command) {
	getGameStateCmd := &cobra.Command{
		Use:   "state",
		Short: "Get the state of a game",
		Long:  "Get the state of a game",
		RunE:  c.Run,
	}

	cli.GameIdFlag(getGameStateCmd.Flags(), &c.GameId)

	parent.AddCommand(getGameStateCmd)
}
