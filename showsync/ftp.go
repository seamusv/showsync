package showsync

import (
	"github.com/secsy/goftp"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
)

type Ftp struct {
	Server   string
	Username string
	Password string
	RootPath string
	Dst      string
}

func (f *Ftp) Sync(path string) error {
	if err := os.MkdirAll(filepath.Join(f.Dst, path), 0755); err != nil {
		return err
	}

	config := goftp.Config{
		User:               f.Username,
		Password:           f.Password,
		ConnectionsPerHost: 10,
		Timeout:            10 * time.Second,
		Logger:             os.Stderr,
	}

	client, err := goftp.DialConfig(config, f.Server)
	if err != nil {
		return err
	}

	err = Walk(client, filepath.Join(f.RootPath, path), func(fullPath string, info os.FileInfo, err error) error {
		if err != nil {
			// no permissions is okay, keep walking
			if err.(goftp.Error).Code() == 550 {
				return nil
			}
			return err
		}

		localPath := strings.TrimPrefix(fullPath, f.RootPath)
		if !info.IsDir() && !strings.HasSuffix(fullPath, ".png") {
			if err := os.MkdirAll(filepath.Join(f.Dst, filepath.Dir(localPath)), 0755); err != nil {
				return err
			}

			remoteFile := fullPath
			localFile := filepath.Join(f.Dst, localPath)

			f, err := os.Create(localFile)
			if err != nil {
				return err
			}
			defer f.Close()

			if err := client.Retrieve(remoteFile, f); err != nil {
				return err
			}
		}

		return nil
	})

	return err
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
