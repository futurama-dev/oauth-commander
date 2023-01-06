package config

import (
	"errors"
	"fmt"
	"github.com/kirsle/configdir"
	"log"
	"os"
	"path/filepath"
)

const (
	FileType     = "yaml"
	FileBaseName = "config"
	FileName     = FileBaseName + "." + FileType
)

func ConfigDir() string {
	localConfig := configdir.LocalConfig("oauth-commander")

	err := configdir.MakePath(localConfig)
	if err != nil {
		panic(err)
	}

	return localConfig
}

func ServerDir() string {
	severDir := configdir.LocalConfig("oauth-commander", "servers")

	err := configdir.MakePath(severDir)
	if err != nil {
		panic(err)
	}

	return severDir
}

func ClientDir(slug string) string {
	clientDir := configdir.LocalConfig("oauth-commander", "servers", slug)

	err := configdir.MakePath(clientDir)
	if err != nil {
		panic(err)
	}

	return clientDir
}

func EnsureConfigFile() {
	configFile := filepath.Join(ConfigDir(), FileName)

	if _, err := os.Stat(configFile); err == nil {
		return
	} else if errors.Is(err, os.ErrNotExist) {
		fmt.Println("Config file not found, creating an empty one:", configFile)
		emptyFile, err := os.Create(configFile)
		if err != nil {
			log.Fatalln(err)
		}
		emptyFile.Close()
	} else {
		log.Fatalln(err)
	}
}
