package server

import (
	"errors"
	"github.com/futurama-dev/oauth-commander/config"
	"github.com/futurama-dev/oauth-commander/discovery"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
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

func (s Server) GetAuthorizationEndpoint() string {
	if endpoint := s.Metadata["authorization_endpoint"]; endpoint == nil {
		return ""
	} else {
		return endpoint.(string)
	}
}

func (s Server) GetSupportedScopes() []string {
	if scopesSupported := s.Metadata["scopes_supported"]; scopesSupported == nil {
		return []string{}
	} else {
		return scopesSupported.([]string)
	}
}

func (s Server) GetSupportedResponseTypes() []string {
	if responseTypesSupported := s.Metadata["response_types_supported"]; responseTypesSupported == nil {
		return []string{}
	} else {
		return responseTypesSupported.([]string)
	}
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

func RemoveBySlug(slug string) error {
	err := remove(slug, config.ServerDir())
	return err
}

func remove(slug string, serverDir string) error {
	pathToFile := filepath.Join(serverDir, slug+".yaml")
	_, err := os.Stat(pathToFile)
	if err == os.ErrNotExist {
		return errors.New("server file not found to remove")
	} else if err == nil {
		err = os.Remove(pathToFile)
		if err != nil {
			return err
		}
		if config.GetSelectedServer() == slug {
			err = config.SetSelectedServer("")
			if err != nil {
				return err
			}
		}
	} else {
		return err
	}

	return nil
}

func Save(server Server) error {
	err := write(server, false, config.ServerDir())

	if err == nil {
		if !config.IsSelectedServer() {
			return config.SetSelectedServer(server.Slug)
		}
	}

	return err
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
