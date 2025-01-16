package cmd

import (
	"github.com/mcoot/crosswordgame-go/internal/cli"
	"github.com/mcoot/crosswordgame-go/internal/client"
	"github.com/spf13/cobra"
)

type HealthCommand struct{}

func (c *HealthCommand) Run(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	cwg := client.GetClient(ctx)

	health, err := cwg.Health()
	if err != nil {
		return err
	}

	return cli.WriteOutput(health)
}

func (c *HealthCommand) Mount(parent *cobra.Command) {
	healthCmd := &cobra.Command{
		Use:   "health",
		Short: "Check the health of the server",
		Long:  "Check the health of the server",
		RunE:  c.Run,
	}
	parent.AddCommand(healthCmd)
}
