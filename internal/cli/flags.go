package cli

import (
	"github.com/spf13/pflag"
)

var (
	FlagServer     string
	FlagOutputMode OutputMode
)

func GlobalFlagServer(flagSet *pflag.FlagSet) {
	flagSet.
		StringVarP(&FlagServer, "server", "s", "http://localhost:8080", "Server URL")
}

func GlobalFlagOutputMode(flagSet *pflag.FlagSet) {
	flagSet.
		VarP(&FlagOutputMode, "output", "o", "Output mode (default: text, allowed: text, json, yaml)")
}

func GameIdFlag(flagSet *pflag.FlagSet, v *string) {
	flagSet.
		StringVarP(v, "game", "g", "", "Game ID")
}

func PlayerIdFlag(flagSet *pflag.FlagSet, v *int) {
	flagSet.
		IntVarP(v, "player", "p", -1, "Player ID")
}
