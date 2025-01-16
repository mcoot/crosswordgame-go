package cmd

import (
	"github.com/mcoot/crosswordgame-go/internal/client"
	"github.com/spf13/cobra"
)

var (
	healthCmd = &cobra.Command{
		Use:   "health",
		Short: "Check the health of the server",
		Long:  "Check the health of the server",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			c := client.GetClient(ctx)

			health, err := c.Health()
			if err != nil {
				return err
			}

			return OutputMode.WriteOutput(health)
		},
	}
)
