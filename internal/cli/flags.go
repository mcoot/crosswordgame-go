package cli

import (
	"github.com/spf13/cobra"
)

var (
	FlagServer     string
	FlagOutputMode OutputMode
)

func GlobalFlagServer(cmd *cobra.Command) {
	cmd.PersistentFlags().
		StringVarP(&FlagServer, "server", "s", "http://localhost:8080", "Server URL")
}

func GlobalFlagOutputMode(cmd *cobra.Command) {
	cmd.PersistentFlags().
		VarP(&FlagOutputMode, "output", "o", "Output mode (default: text, allowed: text, json, yaml)")
}

func GameIdFlag(cmd *cobra.Command, v *string) {
	cmd.Flags().
		StringVarP(v, "game", "g", "", "Game ID")
	err := cmd.MarkFlagRequired("game")
	if err != nil {
		panic(err)
	}
}

func PlayerIdFlag(cmd *cobra.Command, v *int) {
	cmd.Flags().
		IntVarP(v, "player", "p", -1, "Player ID")
	err := cmd.MarkFlagRequired("player")
	if err != nil {
		panic(err)
	}
}
