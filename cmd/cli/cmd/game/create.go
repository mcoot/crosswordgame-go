package game

import "github.com/spf13/cobra"

type CreateGameCommand struct {
	PlayerCount int
}

func (c *CreateGameCommand) Run(cmd *cobra.Command, args []string) error {
	return nil
}

func (c *CreateGameCommand) Mount(parent *cobra.Command) {
	createGameCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new game",
		Long:  "Create a new game",
		RunE:  c.Run,
	}

	createGameCmd.Flags().
		IntVarP(&c.PlayerCount, "players", "p", 2, "Number of players")

	parent.AddCommand(createGameCmd)
}
