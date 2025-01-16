package cli

import "github.com/spf13/cobra"

type Mountable interface {
	Mount(parent *cobra.Command)
}
