package client

import (
	"github.com/futurama-dev/oauth-commander/config"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"
)

type Client struct {
	Slug         string
	Type         string
	CreatedAt    time.Time
	Id           string
	SecretHandle string
}

type Clients []Client

func Load() Clients {
	serverSlug, err := config.CurrentServerSlug()

	if err != nil {
		return Clients{}
	}

	return LoadForServer(serverSlug)
}

func LoadForServer(serverSlug string) Clients {
	return load(serverSlug, config.ServerDir())
}

func load(serverSlug string, serverDir string) Clients {
	clientDir := filepath.Join(serverDir, serverSlug)
	files, err := ioutil.ReadDir(clientDir)

	if err != nil {
		return Clients{}
	}

	var clients Clients

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".yaml") {
			filePath := filepath.Join(clientDir, file.Name())
			clientYaml, err := ioutil.ReadFile(filePath)

			if err != nil {
				log.Println(err)
				continue
			}

			var client Client

			err = yaml.Unmarshal(clientYaml, &client)

			if err != nil {
				log.Println(err)
				continue
			}

			clients = append(clients, client)
		}
	}

	return clients
}
