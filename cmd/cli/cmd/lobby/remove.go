package lobby

import (
	"github.com/mcoot/crosswordgame-go/internal/cli"
	"github.com/mcoot/crosswordgame-go/internal/client"
	lobbytypes "github.com/mcoot/crosswordgame-go/internal/lobby/types"
	playertypes "github.com/mcoot/crosswordgame-go/internal/player/types"
	"github.com/spf13/cobra"
)

type RemovePlayerFromLobbyCommand struct {
	LobbyID  string
	PlayerID string
}

func (c *RemovePlayerFromLobbyCommand) Run(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	cwg := client.GetClient(ctx)

	resp, err := cwg.RemovePlayerFromLobby(lobbytypes.LobbyId(c.LobbyID), playertypes.PlayerId(c.PlayerID))
	if err != nil {
		return err
	}

	return cli.WriteOutput(resp)
}

func (c *RemovePlayerFromLobbyCommand) Mount(parent *cobra.Command) {
	removeCmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove a player from a lobby",
		Long:  "Remove a player from a lobby",
		RunE:  c.Run,
	}

	cli.LobbyIdFlag(removeCmd, &c.LobbyID)
	cli.PlayerIdFlag(removeCmd, &c.PlayerID)

	parent.AddCommand(removeCmd)
}
