package showsync

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

type RSync struct {
	Dst    string
	Src    string
	Passwd string
}

func (r RSync) Sync(path string) bool {
	log.Printf("rsync.Sync: '%s'", path)
	cmd := exec.Command(
		"rsync",
		"-avrms",
		"--password-file",
		r.Passwd,
		"--ignore-existing",
		fmt.Sprintf("%s/%s", r.Src, path),
		r.Dst,
	)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		log.Print(cmd.String())
		log.Print(out.String())
		log.Print(stderr.String())
		log.Print(err)
		return false
	}

	return true
}
