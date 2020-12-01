package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type File struct {
	stagePath     string
	completedPath string
}

func (f File) IsCompleted(path string) bool {
	fullPath := filepath.Join(f.completedPath, path)
	_, err := os.Stat(fullPath)
	return err == nil
}

func (f File) MoveFromStageToDestination(path string) bool {
	log.Printf("file.MoveFromStageToDestination: '%s'", path)
	cmd := exec.Command(
		"mv",
		"-f",
		filepath.Join(f.stagePath, path),
		f.completedPath,
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

func (f File) Unpack(path string) bool {
	log.Printf("file.Unpack: '%s'", path)
	fullPath := filepath.Join(f.stagePath, path)
	info, err := os.Stat(fullPath)
	if err != nil {
		log.Print(err)
		return false
	}
	if info.IsDir() {
		cmd := exec.Command(
			"unrar",
			"x",
			"-o-",
			"*.rar",
		)

		cmd.Dir = fullPath
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok && exiterr.ExitCode() != 10 {
				log.Print(cmd.String())
				log.Print(out.String())
				log.Print(stderr.String())
				log.Print(err)
			}
		}
	}

	return true
}
