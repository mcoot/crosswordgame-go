package lobby

import (
	"github.com/mcoot/crosswordgame-go/internal/cli"
	"github.com/mcoot/crosswordgame-go/internal/client"
	lobbytypes "github.com/mcoot/crosswordgame-go/internal/lobby/types"
	"github.com/spf13/cobra"
)

type DetachGameFromLobbyCommand struct {
	LobbyID string
}

func (c *DetachGameFromLobbyCommand) Run(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	cwg := client.GetClient(ctx)

	resp, err := cwg.DetachGameFromLobby(lobbytypes.LobbyId(c.LobbyID))
	if err != nil {
		return err
	}

	return cli.WriteOutput(resp)
}

func (c *DetachGameFromLobbyCommand) Mount(parent *cobra.Command) {
	detachCmd := &cobra.Command{
		Use:   "detach",
		Short: "Detach a game from a lobby",
		Long:  "Detach a game from a lobby",
		RunE:  c.Run,
	}

	cli.LobbyIdFlag(detachCmd, &c.LobbyID)

	parent.AddCommand(detachCmd)
}
