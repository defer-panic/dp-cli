package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "dp-cli",
	Short: "Toolkit for managing Defer Panic articles and other stuff",
}

func init() {
	rootCmd.AddCommand(articleCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
