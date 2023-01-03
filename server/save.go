package server

import (
	"errors"
	"github.com/futurama-dev/oauth-commander/config"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

func Save(toSave Server) {
	write(toSave, false, config.ServerDir())
}

func Update(toUpdate Server) {
	write(toUpdate, true, config.ServerDir())
}

func write(server Server, overwrite bool, serverDir string) error {
	pathToFile := filepath.Join(serverDir, server.Slug+".yaml")
	_, err := os.Stat(pathToFile)
	if overwrite {
		if err == os.ErrNotExist {
			return errors.New("file not found to update")
		}
	} else {
		if err == nil {
			return errors.New("file found, won't overwrite")
		}
	}

	toWrite, err := yaml.Marshal(&server.Metadata)
	if err != nil {
		return err
	}

	file, err := os.Open(pathToFile)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(string(toWrite))
	if err != nil {
		return err
	}

	return nil
}
