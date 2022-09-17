package cmd

import "github.com/spf13/cobra"

var articleCmd = &cobra.Command{
	Use: "article",
	Short: "Manage articles",
}

func init() {
	articleCmd.AddCommand(articleGenerateCmd)
	articleCmd.AddCommand(articleExportCmd)
}
