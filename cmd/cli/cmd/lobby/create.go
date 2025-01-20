package lobby

import (
	"github.com/mcoot/crosswordgame-go/internal/cli"
	"github.com/mcoot/crosswordgame-go/internal/client"
	"github.com/spf13/cobra"
)

type CreateLobbyCommand struct {
	Name string
}

func (c *CreateLobbyCommand) Run(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	cwg := client.GetClient(ctx)

	resp, err := cwg.CreateLobby(c.Name)
	if err != nil {
		return err
	}

	return cli.WriteOutput(resp)
}

func (c *CreateLobbyCommand) Mount(parent *cobra.Command) {
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new lobby",
		RunE:  c.Run,
	}

	createCmd.Flags().StringVarP(&c.Name, "name", "n", "", "Name of the lobby")

	parent.AddCommand(createCmd)
}
