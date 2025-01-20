package game

import (
	"github.com/mcoot/crosswordgame-go/internal/cli"
	"github.com/mcoot/crosswordgame-go/internal/client"
	playertypes "github.com/mcoot/crosswordgame-go/internal/player/types"
	"github.com/spf13/cobra"
)

type CreateGameCommand struct {
	PlayerIds      []string
	BoardDimension int
}

func (c *CreateGameCommand) Run(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	cwg := client.GetClient(ctx)

	var boardDimension *int
	if c.BoardDimension != 0 {
		boardDimension = &c.BoardDimension
	}

	playerIds := make([]playertypes.PlayerId, len(c.PlayerIds))
	for i, playerId := range c.PlayerIds {
		playerIds[i] = playertypes.PlayerId(playerId)
	}

	game, err := cwg.CreateGame(playerIds, boardDimension)
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
		StringSliceVarP(&c.PlayerIds, "players", "p", []string{},
			"Player IDs in the game, comma-separated",
		)
	_ = createGameCmd.MarkFlagRequired("players")
	createGameCmd.Flags().
		IntVarP(&c.BoardDimension, "dimension", "d", 0, "Board dimension")

	parent.AddCommand(createGameCmd)
}
