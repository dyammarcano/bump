package cmd

import (
	"github.com/spf13/cobra"
	"github.com/treeder/bump/internal/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Bump",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println(version.GetVersionInfo())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
