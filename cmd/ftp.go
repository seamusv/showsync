package cmd

import (
	"fmt"
	"github.com/seamusv/show-sync/showsync"
	"log"
	"net/url"
	"sync"

	"github.com/spf13/cobra"
)

var ftpUrl string

var ftpCmd = &cobra.Command{
	Use:   "ftp",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("URL: %s\n", ftpUrl)
		u, err := url.Parse(ftpUrl)
		if err != nil {
			return err
		}

		if err := showsync.SingleInstance(addr); err != nil {
			return nil
		}

		log.Print("begin")

		password, _ := u.User.Password()
		ftp := &showsync.Ftp{
			Server:   u.Host,
			Username: u.User.Username(),
			Password: password,
			RootPath: u.Path,
			Dst:      stage,
		}

		server := &showsync.Server{
			Api: api,
			Url: serverUrl,
		}

		workQueue := make(chan string, 100)
		completedQueue := make(chan string, 100)
		wg := sync.WaitGroup{}

		for i := 0; i < 5; i++ {
			go func() {
				for {
					select {
					case path := <-workQueue:
						ftp.Sync(path)
						completedQueue <- path
					}
				}
			}()
		}

		file := &showsync.File{
			StagePath:     stage,
			CompletedPath: completed,
		}

		go func() {
			for {
				select {
				case path := <-completedQueue:
					if file.Unpack(path) {
						file.MoveFromStageToDestination(path)
					}
					wg.Done()
				}
			}
		}()

		paths, err := server.GetEntries()
		if err != nil {
			return err
		}

		for _, path := range paths {
			if !file.IsCompleted(path) {
				wg.Add(1)
				workQueue <- path
			}
		}

		wg.Wait()
		log.Print("completed")
		return nil
	},
}

func init() {
	syncCmd.AddCommand(ftpCmd)

	pf := ftpCmd.PersistentFlags()
	pf.StringVar(&ftpUrl, "src", "", "FTP URL")
	_ = ftpCmd.MarkPersistentFlagRequired("src")
}
