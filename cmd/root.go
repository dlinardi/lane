package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lane",
	Short: "Manage parallel development lanes using git worktrees and tmux",
	Long: `lane is a CLI tool for managing parallel development lanes using git worktrees and tmux.

It enables developers to do things like run multiple AI coding sessions simultaneously on different
tickets/branches, switch between them, and clean up when done, without context switching confusion.`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
