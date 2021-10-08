/*
Copyright © 2021 Seamus Venasse <svenasse@gmail.com>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "A brief description of your command",
}

var (
	completed string
	stage     string
	serverUrl string
	api       string
	addr      string
)

func init() {
	rootCmd.AddCommand(syncCmd)

	pf := syncCmd.PersistentFlags()
	pf.StringVar(&completed, "completed", "", "Completed download directory")
	pf.StringVar(&stage, "stage", "", "Staging download directory")
	pf.StringVar(&serverUrl, "url", "", "URL to Sonarr/Radarr server")
	pf.StringVar(&api, "api", "", "API key for Sonarr/Radarr server")
	pf.StringVar(&addr, "addr", ":12345", "Address of instance locking port")
	_ = syncCmd.MarkPersistentFlagRequired("completed")
	_ = syncCmd.MarkPersistentFlagRequired("stage")
	_ = syncCmd.MarkPersistentFlagRequired("url")
	_ = syncCmd.MarkPersistentFlagRequired("api")
}
