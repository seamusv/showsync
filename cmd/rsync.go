package cmd

import (
	"github.com/seamusv/show-sync/showsync"
	"github.com/spf13/cobra"
	"log"
	"sync"
)

var (
	rsyncPasswd string
	rsyncSrc    string
)

var rsyncCmd = &cobra.Command{
	Use:   "rsync",
	Short: "Sync from Sonarr/Radarr and remote server via rsync",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := showsync.SingleInstance(addr); err != nil {
			return nil
		}

		log.Print("begin")

		rsync := &showsync.RSync{
			Dst:    stage,
			Src:    rsyncSrc,
			Passwd: rsyncPasswd,
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
						rsync.Sync(path)
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
	syncCmd.AddCommand(rsyncCmd)

	pf := rsyncCmd.PersistentFlags()
	pf.StringVar(&rsyncPasswd, "passwd", "", "Full path to rSync password file")
	pf.StringVar(&rsyncSrc, "src", "", "Remote source prefix")
	_ = rsyncCmd.MarkPersistentFlagRequired("passwd")
	_ = rsyncCmd.MarkPersistentFlagRequired("src")
}
