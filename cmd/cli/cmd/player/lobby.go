package player

import (
	"github.com/mcoot/crosswordgame-go/internal/cli"
	"github.com/mcoot/crosswordgame-go/internal/client"
	playertypes "github.com/mcoot/crosswordgame-go/internal/player/types"
	"github.com/spf13/cobra"
)

type GetLobbyForPlayerCommand struct {
	PlayerID string
}

func (c *GetLobbyForPlayerCommand) Run(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	cwg := client.GetClient(ctx)

	resp, err := cwg.GetLobbyForPlayer(playertypes.PlayerId(c.PlayerID))
	if err != nil {
		return err
	}

	return cli.WriteOutput(resp)
}

func (c *GetLobbyForPlayerCommand) Mount(parent *cobra.Command) {
	getLobbyCmd := &cobra.Command{
		Use:   "lobby",
		Short: "Get the lobby a player is currently in",
		Long:  "Get the lobby a player is currently in",
		RunE:  c.Run,
	}

	cli.PlayerIdFlag(getLobbyCmd, &c.PlayerID)

	parent.AddCommand(getLobbyCmd)
}
