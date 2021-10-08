package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Build string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version of application",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Build: %s\n", Build)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
