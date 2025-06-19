package commands

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "cli",
		Short: "CLI for AI Task Marketplace",
	}

	rootCmd.AddCommand(submitCmd)
	rootCmd.AddCommand(statusCmd)

	return rootCmd
}
