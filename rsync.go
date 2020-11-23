package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

type RSync struct {
	dst    string
	src    string
	passwd string
}

func (r RSync) Sync(path string) bool {
	log.Printf("rsync.Sync: '%s'", path)
	cmd := exec.Command(
		"rsync",
		"-avrms",
		"--password-file",
		r.passwd,
		"--ignore-existing",
		fmt.Sprintf("%s/%s", r.src, path),
		r.dst,
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
