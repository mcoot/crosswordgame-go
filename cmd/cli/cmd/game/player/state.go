package player

import (
	"github.com/mcoot/crosswordgame-go/internal/cli"
	"github.com/spf13/cobra"
)

type GetPlayerStateCommand struct {
	GameId   string
	PlayerId int
}

func (c *GetPlayerStateCommand) Run(cmd *cobra.Command, args []string) error {
	return nil
}

func (c *GetPlayerStateCommand) Mount(parent *cobra.Command) {
	getPlayerStateCmd := &cobra.Command{
		Use:   "state",
		Short: "Get the state of a player",
		Long:  "Get the state of a player",
		RunE:  c.Run,
	}

	cli.GameIdFlag(getPlayerStateCmd.Flags(), &c.GameId)
	cli.PlayerIdFlag(getPlayerStateCmd.Flags(), &c.PlayerId)

	parent.AddCommand(getPlayerStateCmd)
}
