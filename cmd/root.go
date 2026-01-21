package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version information set via ldflags at build time.
// Falls back to "dev" when not set.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// Verbose flag for debug logging (used by subcommands).
var Verbose bool

var rootCmd = &cobra.Command{
	Use:   "lane",
	Short: "Manage parallel development lanes using git worktrees and tmux",
	Long: `lane is a CLI tool for managing parallel development lanes using git worktrees and tmux.

It enables developers to do things like run multiple AI coding sessions simultaneously on different
tickets/branches, switch between them, and clean up when done, without context switching confusion.`,
	Version:       version,
	SilenceUsage:  true,
	SilenceErrors: true,
	Args:          cobra.NoArgs,
	// Show help when no subcommand is provided
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.SetVersionTemplate("lane version {{.Version}}\n")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "V", false, "enable verbose output for debugging")
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// SetVersionInfo allows setting version information programmatically (useful for testing).
func SetVersionInfo(v, c, d string) {
	version = v
	commit = c
	date = d
	rootCmd.Version = v
}

// GetVersionInfo returns the current version information.
func GetVersionInfo() (ver, com, dat string) {
	return version, commit, date
}

// GetRootCmd returns the root command for testing purposes.
func GetRootCmd() *cobra.Command {
	return rootCmd
}

// NewRootCmd creates a fresh root command instance for testing.
// This avoids state pollution between tests.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lane",
		Short: "Manage parallel development lanes using git worktrees and tmux",
		Long: `lane is a CLI tool for managing parallel development lanes using git worktrees and tmux.

It enables developers to do things like run multiple AI coding sessions simultaneously on different
tickets/branches, switch between them, and clean up when done, without context switching confusion.`,
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.SetVersionTemplate("lane version {{.Version}}\n")
	cmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "V", false, "enable verbose output for debugging")
	return cmd
}

// ResetRootCmd resets the root command state (useful for testing).
func ResetRootCmd() {
	rootCmd.SetArgs([]string{})
	rootCmd.SetOut(nil)
	rootCmd.SetErr(nil)
}

// VersionFull returns the full version string including commit and date.
func VersionFull() string {
	return fmt.Sprintf("lane version %s (commit: %s, built: %s)", version, commit, date)
}
