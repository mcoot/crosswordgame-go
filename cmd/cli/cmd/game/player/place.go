package player

import (
	"github.com/mcoot/crosswordgame-go/internal/cli"
	"github.com/mcoot/crosswordgame-go/internal/client"
	"github.com/mcoot/crosswordgame-go/internal/game/types"
	playertypes "github.com/mcoot/crosswordgame-go/internal/player/types"
	"github.com/spf13/cobra"
)

type PlayerPlaceCommand struct {
	GameId   string
	PlayerId string
	Row      int
	Column   int
}

func (c *PlayerPlaceCommand) Run(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	cwg := client.GetClient(ctx)

	resp, err := cwg.SubmitPlacement(types.GameId(c.GameId), playertypes.PlayerId(c.PlayerId), c.Row, c.Column)
	if err != nil {
		return err
	}

	return cli.WriteOutput(resp)
}

func (c *PlayerPlaceCommand) Mount(parent *cobra.Command) {
	playerPlaceCmd := &cobra.Command{
		Use:   "place",
		Short: "Place a letter",
		Long:  "Place a letter for a player",
		RunE:  c.Run,
	}

	cli.GameIdFlag(playerPlaceCmd, &c.GameId)
	cli.PlayerIdFlag(playerPlaceCmd, &c.PlayerId)

	playerPlaceCmd.Flags().IntVarP(&c.Row, "row", "r", 0, "Row")
	_ = playerPlaceCmd.MarkFlagRequired("row")
	playerPlaceCmd.Flags().IntVarP(&c.Column, "column", "c", 0, "Column")
	_ = playerPlaceCmd.MarkFlagRequired("column")

	parent.AddCommand(playerPlaceCmd)
}
