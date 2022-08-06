package showsync

import (
	"bytes"
	"github.com/rs/zerolog/log"
	"os"
	"os/exec"
	"path/filepath"
)

func LocalPathExists(localPath string, path string) bool {
	_, err := os.Stat(filepath.Join(localPath, path))
	return err == nil
}

func Move(fromPath, toPath, path string) error {
	cmd := exec.Command(
		"mv",
		"-f",
		filepath.Join(fromPath, path),
		filepath.Join(toPath, path),
	)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		log.Err(err).Str("cmd", cmd.String()).Str("stdout", out.String()).Str("stderr", stderr.String()).Msg("failed to move torrent")
		return err
	}

	return nil
}

func Unpack(localPath, path string) error {
	fullPath := filepath.Join(localPath, path)
	info, err := os.Stat(fullPath)
	if err != nil {
		return err
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
				log.Err(err).Str("cmd", cmd.String()).Str("stdout", out.String()).Str("stderr", stderr.String()).Msg("failed to unpack torrent")
				return err
			}
		}
	}

	return nil
}
