package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kak-tus/odiag-split/opendiag"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	if len(os.Args) != 2 {
		return
	}

	if err := process(os.Args[1]); err != nil {
		logger.Error().Err(err).Msg("process failed")
		return
	}
}

func process(dir string) error {
	dirs, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range dirs {
		if entry.IsDir() {
			continue
		}

		if !opendiag.SupportedFileName(entry.Name()) {
			continue
		}

		parsed, err := opendiag.DateFromFileName(entry.Name())
		if err != nil {
			return err
		}

		toRead := filepath.Join(dir, entry.Name())

		content, err := os.ReadFile(toRead)
		if err != nil {
			return err
		}

		olog, err := opendiag.Decode(parsed, string(content))
		if err != nil {
			return err
		}

		if !olog.NeedSplit() {
			continue
		}

		splitted := olog.Split()

		for _, newLog := range splitted {
			newName, encoded := newLog.Encode()

			newName = filepath.Join(dir, newName)
			if _, err := os.Stat(newName); err == nil {
				return fmt.Errorf("strange, but file '%s' exists", newName)
			}

			if err := os.WriteFile(newName, []byte(encoded), 0644); err != nil {
				return err
			}
		}

		toBackup := toRead + ".backup"
		if err := os.Rename(toRead, toBackup); err != nil {
			return err
		}
	}

	return nil
}
