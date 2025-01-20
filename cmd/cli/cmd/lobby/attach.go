package lobby

import (
	"github.com/mcoot/crosswordgame-go/internal/cli"
	"github.com/mcoot/crosswordgame-go/internal/client"
	gametypes "github.com/mcoot/crosswordgame-go/internal/game/types"
	lobbytypes "github.com/mcoot/crosswordgame-go/internal/lobby/types"
	"github.com/spf13/cobra"
)

type AttachGameToLobbyCommand struct {
	LobbyID string
	GameID  string
}

func (c *AttachGameToLobbyCommand) Run(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	cwg := client.GetClient(ctx)

	resp, err := cwg.AttachGameToLobby(lobbytypes.LobbyId(c.LobbyID), gametypes.GameId(c.GameID))
	if err != nil {
		return err
	}

	return cli.WriteOutput(resp)
}

func (c *AttachGameToLobbyCommand) Mount(parent *cobra.Command) {
	attachCmd := &cobra.Command{
		Use:   "attach",
		Short: "Attach a game to a lobby",
		Long:  "Attach a game to a lobby",
		RunE:  c.Run,
	}

	cli.LobbyIdFlag(attachCmd, &c.LobbyID)
	cli.GameIdFlag(attachCmd, &c.GameID)

	parent.AddCommand(attachCmd)
}
