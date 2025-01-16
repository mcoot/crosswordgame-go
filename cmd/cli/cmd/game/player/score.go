package player

import (
	"github.com/mcoot/crosswordgame-go/internal/cli"
	"github.com/spf13/cobra"
)

type GetPlayerScoreCommand struct {
	GameId   string
	PlayerId int
}

func (c *GetPlayerScoreCommand) Run(cmd *cobra.Command, args []string) error {
	return nil
}

func (c *GetPlayerScoreCommand) Mount(parent *cobra.Command) {
	getPlayerScoreCmd := &cobra.Command{
		Use:   "score",
		Short: "Get player score",
		Long:  "Get the score of a player",
		RunE:  c.Run,
	}

	cli.GameIdFlag(getPlayerScoreCmd.Flags(), &c.GameId)
	cli.PlayerIdFlag(getPlayerScoreCmd.Flags(), &c.PlayerId)

	parent.AddCommand(getPlayerScoreCmd)
}
