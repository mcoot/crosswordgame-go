package player

import (
	"github.com/mcoot/crosswordgame-go/internal/cli"
	"github.com/mcoot/crosswordgame-go/internal/client"
	"github.com/mcoot/crosswordgame-go/internal/game/types"
	"github.com/spf13/cobra"
)

type GetPlayerStateCommand struct {
	GameId   string
	PlayerId int
}

func (c *GetPlayerStateCommand) Run(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	cwg := client.GetClient(ctx)

	state, err := cwg.GetPlayerState(types.GameId(c.GameId), c.PlayerId)
	if err != nil {
		return err
	}

	return cli.WriteOutput(state)
}

func (c *GetPlayerStateCommand) Mount(parent *cobra.Command) {
	getPlayerStateCmd := &cobra.Command{
		Use:   "state",
		Short: "Get the state of a player",
		Long:  "Get the state of a player",
		RunE:  c.Run,
	}

	cli.GameIdFlag(getPlayerStateCmd, &c.GameId)
	cli.PlayerIdFlag(getPlayerStateCmd, &c.PlayerId)

	parent.AddCommand(getPlayerStateCmd)
}
