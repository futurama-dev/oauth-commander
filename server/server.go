package server

import (
	"github.com/futurama-dev/oauth-commander/config"
	"github.com/futurama-dev/oauth-commander/discovery"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Server struct {
	Slug      string
	Type      string
	CreatedAt time.Time `yaml:"created_at"`
	Metadata  map[string]any
}

type Servers []Server

func Load() Servers {
	return load(config.ServerDir())
}

func load(serverDir string) Servers {
	files, err := ioutil.ReadDir(serverDir)

	if err != nil {
		log.Fatalln(err)
	}

	var servers []Server

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".yaml") {
			filePath := filepath.Join(serverDir, file.Name())
			serverYaml, err := ioutil.ReadFile(filePath)

			if err != nil {
				log.Println(err)
				continue
			}

			var server Server

			err = yaml.Unmarshal(serverYaml, &server)

			if err != nil {
				log.Println(err)
				continue
			}

			servers = append(servers, server)
		}
	}

	return servers
}

func (servers Servers) FindBySlug(slug string) (Server, bool) {
	for _, server := range servers {
		if server.Slug == slug {
			return server, true
		}
	}

	return Server{}, false
}

func (servers Servers) FindByIssuer(slug string) (Server, bool) {
	for _, server := range servers {
		if server.Metadata["issuer"] == slug {
			return server, true
		}
	}

	return Server{}, false
}

func IssuerToSlug(issuer string) (string, error) {
	_, err := discovery.ValidateIssuer(issuer)

	if err != nil {
		return "", err
	}

	issuer = issuer[8:]

	var re = regexp.MustCompile(`[^a-z0-9_]+`)

	return strings.TrimRight(re.ReplaceAllString(issuer, "_"), "_"), nil
}
