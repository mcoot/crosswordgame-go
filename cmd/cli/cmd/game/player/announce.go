package player

import (
	"github.com/mcoot/crosswordgame-go/internal/cli"
	"github.com/mcoot/crosswordgame-go/internal/client"
	"github.com/mcoot/crosswordgame-go/internal/game/types"
	"github.com/spf13/cobra"
)

type PlayerAnnounceCommand struct {
	GameId   string
	PlayerId int
	Letter   cli.LetterValue
}

func (c *PlayerAnnounceCommand) Run(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	cwg := client.GetClient(ctx)

	resp, err := cwg.SubmitAnnouncement(types.GameId(c.GameId), c.PlayerId, string(c.Letter))
	if err != nil {
		return err
	}

	return cli.WriteOutput(resp)
}

func (c *PlayerAnnounceCommand) Mount(parent *cobra.Command) {
	playerAnnounceCmd := &cobra.Command{
		Use:   "announce",
		Short: "Announce a letter",
		Long:  "Announce a letter for a player",
		RunE:  c.Run,
	}

	cli.GameIdFlag(playerAnnounceCmd, &c.GameId)
	cli.PlayerIdFlag(playerAnnounceCmd, &c.PlayerId)
	cli.LetterFlag(playerAnnounceCmd, &c.Letter)

	parent.AddCommand(playerAnnounceCmd)
}
