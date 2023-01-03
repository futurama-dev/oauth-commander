package server

import (
	"errors"
	"github.com/futurama-dev/oauth-commander/config"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

func Save(server Server) error {
	return write(server, false, config.ServerDir())
}

func Update(server Server) error {
	return write(server, true, config.ServerDir())
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

	toWrite, err := yaml.Marshal(&server)
	if err != nil {
		return err
	}

	return os.WriteFile(pathToFile, toWrite, 0644)
}
