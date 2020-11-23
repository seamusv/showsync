package main

import (
	"flag"
	"log"
	"os"
	"sync"
)

var (
	rsyncPasswordFile = flag.String("passwd", "", "rSync password File")
	rsyncSource       = flag.String("src", "", "rSync source Prefix")
	stagePath         = flag.String("stage", "", "rSync destination")
	completedPath     = flag.String("completed", "", "Completed destination")
	url               = flag.String("url", "", "URL of Sonarr/Radarr server")
	api               = flag.String("api", "", "API key for server access")
	addr              = flag.String("addr", ":12345", "instance port for concurrency")
)

func main() {
	if err := SingleInstance(*addr); err != nil {
		panic(err)
	}

	flag.Parse()
	if anyNotSet(*rsyncPasswordFile, *rsyncSource, *stagePath, *completedPath, *url, *api) {
		flag.Usage()
		os.Exit(-1)
	}

	log.Print("begin")

	rsync := &RSync{
		dst:    *stagePath,
		src:    *rsyncSource,
		passwd: *rsyncPasswordFile,
	}

	server := &Server{
		api: *api,
		url: *url,
	}

	file := &File{
		stagePath:     *stagePath,
		completedPath: *completedPath,
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
		panic(err)
	}

	for _, path := range paths {
		if !file.IsCompleted(path) {
			wg.Add(1)
			workQueue <- path
		}
	}

	wg.Wait()
	log.Print("completed")
}

func anyNotSet(opts ...string) bool {
	for _, opt := range opts {
		if len(opt) == 0 {
			return true
		}
	}
	return false
}
