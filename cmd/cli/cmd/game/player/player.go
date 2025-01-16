package player

import "github.com/spf13/cobra"

type PlayerCommand struct{}

func (c *PlayerCommand) Mount(parent *cobra.Command) {
	playerCmd := &cobra.Command{
		Use:   "player",
		Short: "Player commands",
		Long:  "Commands for interacting with a player",
	}

	(&GetPlayerStateCommand{}).Mount(playerCmd)
	(&GetPlayerScoreCommand{}).Mount(playerCmd)

	parent.AddCommand(playerCmd)
}
