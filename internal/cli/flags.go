package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
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

type LetterValue string

func (l *LetterValue) Set(value string) error {
	*l = LetterValue(strings.ToUpper(value))
	if len(*l) != 1 {
		return fmt.Errorf("letter value must be single character, got: %s", value)
	}
	if (*l)[0] < 'A' || (*l)[0] > 'Z' {
		return fmt.Errorf("letter value must be a letter, got: %s", value)
	}
	return nil
}

func (l *LetterValue) String() string {
	return string(*l)
}

func (l *LetterValue) Type() string {
	return "letterValue"
}

func LetterFlag(cmd *cobra.Command, v *LetterValue) {
	cmd.Flags().
		VarP(v, "letter", "l", "Letter")
	err := cmd.MarkFlagRequired("letter")
	if err != nil {
		panic(err)
	}
}
