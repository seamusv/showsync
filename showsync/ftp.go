package showsync

import (
	"crypto/tls"
	"github.com/secsy/goftp"
	"golang.org/x/sync/errgroup"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync/atomic"
	"time"
)

func MakeClient(serverUrl *url.URL) (*goftp.Client, error) {
	config := goftp.Config{
		User:               serverUrl.User.Username(),
		ConnectionsPerHost: 10,
		Timeout:            10 * time.Second,
		Logger:             nil,
	}
	config.Password, _ = serverUrl.User.Password()

	if serverUrl.Scheme == "ftps" {
		config.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
		config.TLSMode = goftp.TLSExplicit
	}

	return goftp.DialConfig(config, serverUrl.Host)
}

func PrepareFtpQueue(serverUrl *url.URL, paths []string) ([]FtpFileInfo, error) {
	var client *goftp.Client
	{
		var err error
		client, err = MakeClient(serverUrl)
		if err != nil {
			return nil, err
		}
	}

	transferQueue := make([]FtpFileInfo, 0)
	{
		err := Walk(client, serverUrl.Path, func(fullPath string, info os.FileInfo, err error) error {
			if err != nil {
				// no permissions, keep walking
				if err.(goftp.Error).Code() == 550 {
					return nil
				}
				return err
			}

			localPath := strings.TrimPrefix(fullPath, serverUrl.Path)
			localPath = strings.TrimPrefix(localPath, "/")

			if info.IsDir() && !inPathsPrefix(paths, localPath) {
				return filepath.SkipDir
			}

			if inPathsPrefix(paths, localPath) {
				if info.IsDir() {
					transferQueue = append(transferQueue, FtpFileInfo{Path: localPath, IsDir: true})
				} else {
					transferQueue = append(transferQueue, FtpFileInfo{Path: localPath, Size: info.Size()})
				}
			}

			return nil
		})
		if err != nil {
			return nil, err
		}

		sort.Slice(transferQueue, func(i, j int) bool {
			return transferQueue[i].Path < transferQueue[j].Path
		})
	}

	return transferQueue, nil
}

func ProcessQueue(serverUrl *url.URL, localPath string, transferQueue []FtpFileInfo) error {
	var client *goftp.Client
	{
		var err error
		client, err = MakeClient(serverUrl)
		if err != nil {
			return err
		}
	}

	queue := make(chan FtpFileInfo, 1)
	g := &errgroup.Group{}
	for i := 0; i < 5; i++ {
		g.Go(func() error {
			for {
				select {
				case item, ok := <-queue:
					if !ok {
						return nil
					}
					localFile := filepath.Join(localPath, item.Path)
					if item.IsDir {
						if err := os.MkdirAll(localFile, 0755); err != nil {
							return err
						}
					} else {
						info, err := os.Stat(localFile)
						if err != nil || info.Size() != item.Size {
							if err := Retrieve(client, filepath.Join(serverUrl.Path, item.Path), localFile); err != nil {
								return err
							}
						}
					}
				}
			}
		})
	}

	for _, info := range transferQueue {
		queue <- info
	}
	close(queue)

	return g.Wait()
}

func Walk(client *goftp.Client, root string, walkFn filepath.WalkFunc) (ret error) {
	dirsToCheck := make(chan string, 100)

	var workCount int32 = 1
	dirsToCheck <- root

	for dir := range dirsToCheck {
		go func(dir string) {
			files, err := client.ReadDir(dir)

			if err != nil {
				if err = walkFn(dir, nil, err); err != nil && err != filepath.SkipDir {
					ret = err
					close(dirsToCheck)
					return
				}
			}

			for _, file := range files {
				if err = walkFn(path.Join(dir, file.Name()), file, nil); err != nil {
					if file.IsDir() && err == filepath.SkipDir {
						continue
					}
					ret = err
					close(dirsToCheck)
					return
				}

				if file.IsDir() {
					atomic.AddInt32(&workCount, 1)
					dirsToCheck <- path.Join(dir, file.Name())
				}
			}

			atomic.AddInt32(&workCount, -1)
			if workCount == 0 {
				close(dirsToCheck)
			}
		}(dir)
	}

	return ret
}

func Retrieve(client *goftp.Client, remoteFile string, localFile string) error {
	f, err := os.Create(localFile)
	if err != nil {
		return err
	}
	defer f.Close()

	return client.Retrieve(remoteFile, f)
}

func inPathsPrefix(paths []string, path string) bool {
	for _, p := range paths {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

type FtpFileInfo struct {
	Path  string
	IsDir bool
	Size  int64
}
