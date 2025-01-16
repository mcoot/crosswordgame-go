package game

import (
	"github.com/mcoot/crosswordgame-go/internal/cli"
	"github.com/mcoot/crosswordgame-go/internal/client"
	"github.com/mcoot/crosswordgame-go/internal/game/types"
	"github.com/spf13/cobra"
)

type GetGameStateCommand struct {
	GameId string
}

func (c *GetGameStateCommand) Run(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	cwg := client.GetClient(ctx)

	state, err := cwg.GetGameState(types.GameId(c.GameId))
	if err != nil {
		return err
	}

	return cli.WriteOutput(state)
}

func (c *GetGameStateCommand) Mount(parent *cobra.Command) {
	getGameStateCmd := &cobra.Command{
		Use:   "state",
		Short: "Get the state of a game",
		Long:  "Get the state of a game",
		RunE:  c.Run,
	}

	cli.GameIdFlag(getGameStateCmd, &c.GameId)

	parent.AddCommand(getGameStateCmd)
}
