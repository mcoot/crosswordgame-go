package lobby

import "github.com/spf13/cobra"

type LobbyCommand struct{}

func (c *LobbyCommand) Mount(parent *cobra.Command) {
	lobbyCmd := &cobra.Command{
		Use:   "lobby",
		Short: "Lobby commands",
		Long:  "Commands for interacting with a lobby",
	}

	(&CreateLobbyCommand{}).Mount(lobbyCmd)
	(&GetLobbyStateCommand{}).Mount(lobbyCmd)
	(&JoinLobbyCommand{}).Mount(lobbyCmd)
	(&RemovePlayerFromLobbyCommand{}).Mount(lobbyCmd)
	(&AttachGameToLobbyCommand{}).Mount(lobbyCmd)
	(&DetachGameFromLobbyCommand{}).Mount(lobbyCmd)

	parent.AddCommand(lobbyCmd)
}
