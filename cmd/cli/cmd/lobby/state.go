package lobby

import (
	"github.com/mcoot/crosswordgame-go/internal/cli"
	"github.com/mcoot/crosswordgame-go/internal/client"
	lobbytypes "github.com/mcoot/crosswordgame-go/internal/lobby/types"
	"github.com/spf13/cobra"
)

type GetLobbyStateCommand struct {
	LobbyID string
}

func (c *GetLobbyStateCommand) Run(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	cwg := client.GetClient(ctx)

	resp, err := cwg.GetLobbyState(lobbytypes.LobbyId(c.LobbyID))
	if err != nil {
		return err
	}

	return cli.WriteOutput(resp)
}

func (c *GetLobbyStateCommand) Mount(parent *cobra.Command) {
	getLobbyStateCmd := &cobra.Command{
		Use:   "state",
		Short: "Get the state of a lobby",
		RunE:  c.Run,
	}

	cli.LobbyIdFlag(getLobbyStateCmd, &c.LobbyID)

	parent.AddCommand(getLobbyStateCmd)
}
