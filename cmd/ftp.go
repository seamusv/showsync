package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/seamusv/show-sync/showsync"
	"github.com/spf13/cobra"
	"net/url"
)

var ftpUrl string

var ftpCmd = &cobra.Command{
	Use:   "ftp",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := showsync.SingleInstance(addr); err != nil {
			return nil
		}

		completedTorrents := make([]string, 0)
		{
			darrServerUrl, err := url.Parse(serverUrl)
			if err != nil {
				return err
			}

			torrents, err := showsync.GetCompletedTorrents(darrServerUrl)
			if err != nil {
				log.Err(err).Msg("failed to get completed torrents")
				return nil
			}

			for _, torrent := range torrents {
				if !showsync.LocalPathExists(completed, torrent) {
					completedTorrents = append(completedTorrents, torrent)
				}
			}
		}

		ftpServerUrl, err := url.Parse(ftpUrl)
		if err != nil {
			return err
		}

		transferQueue, err := showsync.PrepareFtpQueue(ftpServerUrl, completedTorrents)
		if err != nil {
			log.Err(err).Msg("error preparing ftp queue")
			return nil
		}

		err = showsync.ProcessQueue(ftpServerUrl, stage, transferQueue)
		if err != nil {
			log.Err(err).Msg("error processing queue")
			return nil
		}

		for _, torrent := range completedTorrents {
			if err := showsync.Unpack(stage, torrent); err == nil {
				if err := showsync.Move(stage, completed, torrent); err != nil {
					return nil
				}
			}

		}

		return nil
	},
}

func init() {
	syncCmd.AddCommand(ftpCmd)

	pf := ftpCmd.PersistentFlags()
	pf.StringVar(&ftpUrl, "src", "", "FTP URL")
	_ = ftpCmd.MarkPersistentFlagRequired("src")
}
