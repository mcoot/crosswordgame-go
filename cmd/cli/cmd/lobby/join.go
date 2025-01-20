package lobby

import (
	"github.com/mcoot/crosswordgame-go/internal/cli"
	"github.com/mcoot/crosswordgame-go/internal/client"
	lobbytypes "github.com/mcoot/crosswordgame-go/internal/lobby/types"
	playertypes "github.com/mcoot/crosswordgame-go/internal/player/types"
	"github.com/spf13/cobra"
)

type JoinLobbyCommand struct {
	LobbyID  string
	PlayerID string
}

func (c *JoinLobbyCommand) Run(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	cwg := client.GetClient(ctx)

	resp, err := cwg.JoinLobby(lobbytypes.LobbyId(c.LobbyID), playertypes.PlayerId(c.PlayerID))
	if err != nil {
		return err
	}

	return cli.WriteOutput(resp)
}

func (c *JoinLobbyCommand) Mount(parent *cobra.Command) {
	joinCmd := &cobra.Command{
		Use:   "join",
		Short: "Join a player into a lobby",
		Long:  "Join a player into a lobby",
		RunE:  c.Run,
	}

	cli.LobbyIdFlag(joinCmd, &c.LobbyID)
	cli.PlayerIdFlag(joinCmd, &c.PlayerID)

	parent.AddCommand(joinCmd)
}
