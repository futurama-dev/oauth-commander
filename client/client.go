package client

import (
	"github.com/futurama-dev/oauth-commander/config"
	"github.com/zalando/go-keyring"
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
	serverSlug := config.GetCurrentServer()

	if len(serverSlug) == 0 {
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

func (clients Clients) FindBySlug(slug string) (Client, bool) {
	for _, client := range clients {
		if client.Slug == slug {
			return client, true
		}
	}

	return Client{}, false
}

func (clients Clients) FindById(clientId string) (Client, bool) {
	for _, client := range clients {
		if client.Id == clientId {
			return client, true
		}
	}

	return Client{}, false
}

func (client Client) Secret() (string, error) {
	return keyring.Get("oauth-commander", client.SecretHandle)
}

func (client Client) SetSecret(secret string) error {
	return keyring.Set("oauth-commander", client.SecretHandle, secret)
}
