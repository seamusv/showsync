package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "show-sync",
	Short: "Torrent synchronisation tool",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
