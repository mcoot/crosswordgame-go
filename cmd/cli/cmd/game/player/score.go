package player

import (
	"github.com/mcoot/crosswordgame-go/internal/cli"
	"github.com/mcoot/crosswordgame-go/internal/client"
	"github.com/mcoot/crosswordgame-go/internal/game/types"
	"github.com/spf13/cobra"
)

type GetPlayerScoreCommand struct {
	GameId   string
	PlayerId int
}

func (c *GetPlayerScoreCommand) Run(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	cwg := client.GetClient(ctx)

	score, err := cwg.GetPlayerScore(types.GameId(c.GameId), c.PlayerId)
	if err != nil {
		return err
	}

	return cli.WriteOutput(score)
}

func (c *GetPlayerScoreCommand) Mount(parent *cobra.Command) {
	getPlayerScoreCmd := &cobra.Command{
		Use:   "score",
		Short: "Get player score",
		Long:  "Get the score of a player",
		RunE:  c.Run,
	}

	cli.GameIdFlag(getPlayerScoreCmd, &c.GameId)
	cli.PlayerIdFlag(getPlayerScoreCmd, &c.PlayerId)

	parent.AddCommand(getPlayerScoreCmd)
}
