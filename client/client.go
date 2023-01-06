package client

import (
	"github.com/futurama-dev/oauth-commander/config"
	"github.com/spf13/viper"
	"github.com/zalando/go-keyring"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	Slug         string
	ServerSlug   string
	CreatedAt    time.Time
	Id           string
	SecretHandle string
	RedirectURIs []string
}

type Clients []Client

func Load() Clients {
	serverSlug := viper.GetString(config.SelectedServerSlug)

	if len(serverSlug) == 0 {
		return Clients{}
	}

	return LoadForServer(serverSlug)
}

func LoadForServer(serverSlug string) Clients {
	return load(config.ClientDir(serverSlug))
}

func load(clientDir string) Clients {
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

func FindBySlug(slug string) (Client, bool) {
	return Load().FindBySlug(slug)
}

func FindById(clientId string) (Client, bool) {
	return Load().FindById(clientId)
}

func (clients Clients) NextSlug() string {
	i := 0
	ok := true
	slug := ""
	for ok {
		i += 1
		slug = "client_" + strconv.Itoa(i)
		_, ok = clients.FindBySlug(slug)
	}

	return slug
}
func (client Client) Secret() (string, error) {
	return keyring.Get("oauth-commander", client.SecretHandle)
}

func (client Client) SetSecret(secret string) error {
	return keyring.Set("oauth-commander", client.SecretHandle, secret)
}

func NextSlug() string {
	return Load().NextSlug()
}

func Save(client Client) error {
	pathToFile := filepath.Join(config.ClientDir(client.ServerSlug), client.Slug+".yaml")

	toWrite, err := yaml.Marshal(&client)
	if err != nil {
		return err
	}

	return os.WriteFile(pathToFile, toWrite, 0644)
}
