package server

import (
	"github.com/futurama-dev/oauth-commander/config"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type Server struct {
	Slug      string
	Type      string
	CreatedAt time.Time
	Metadata  map[string]any
}

func Load() []Server {
	files, err := ioutil.ReadDir(config.ServerDir())

	if err != nil {
		log.Fatalln(err)
	}

	var servers []Server

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".yaml") {
			serverYaml, err := ioutil.ReadFile(file.Name())

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
