package game

import (
	"github.com/mcoot/crosswordgame-go/cmd/cli/cmd/game/player"
	"github.com/spf13/cobra"
)

type GameCommand struct{}

func (c *GameCommand) Mount(parent *cobra.Command) {
	gameCmd := &cobra.Command{
		Use:   "game",
		Short: "Game commands",
		Long:  "Commands for interacting with a game",
	}

	(&CreateGameCommand{}).Mount(gameCmd)
	(&GetGameStateCommand{}).Mount(gameCmd)
	(&player.PlayerCommand{}).Mount(gameCmd)

	parent.AddCommand(gameCmd)
}
