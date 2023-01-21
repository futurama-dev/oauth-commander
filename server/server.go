package server

import (
	"errors"
	"github.com/futurama-dev/oauth-commander/config"
	"github.com/futurama-dev/oauth-commander/discovery"
	pkce "github.com/nirasan/go-oauth-pkce-code-verifier"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
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

func (s Server) GetTokenEndpoint() string {
	if endpoint := s.Metadata["token_endpoint"]; endpoint == nil {
		return ""
	} else {
		return endpoint.(string)
	}
}

func (s Server) GetSupportedScopes() []string {
	scopesSupported, err := extractStringSlice(s.Metadata["scopes_supported"])
	if err != nil {
		// TODO log or maybe error
		return []string{}
	}

	return scopesSupported
}

func (s Server) GetSupportedResponseTypes() []string {
	typesSupported, err := extractStringSlice(s.Metadata["response_types_supported"])
	if err != nil {
		// TODO log or maybe error
		return []string{}
	}

	return typesSupported
}

func (s Server) GetSupportedCodeChallengeMethods() []string {
	methodsSupported, err := extractStringSlice(s.Metadata["code_challenge_methods_supported"])
	if err != nil {
		// TODO log or maybe error
		return []string{}
	}

	return methodsSupported
}

func (s Server) IsPkceSupported() bool {
	return len(s.GetSupportedCodeChallengeMethods()) > 0
}

func (s Server) IsPkceMethodS256Supported() bool {
	if s.IsPkceSupported() {
		for _, method := range s.GetSupportedCodeChallengeMethods() {
			if method == "S256" {
				return true
			}
		}
	}

	return false
}

func (s Server) GenerateCodeVerifier() (string, string, string, error) {
	if !s.IsPkceSupported() {
		return "", "", "", errors.New("PKCE not supported")
	}

	var codeChallengeMethod, codeVerifier, codeChallenge string

	if s.IsPkceMethodS256Supported() {
		codeChallengeMethod = "S256"
	} else {
		codeChallengeMethod = s.GetSupportedCodeChallengeMethods()[0]
	}

	verifier, err := pkce.CreateCodeVerifier()
	if err != nil {
		return "", "", "", err
	}

	switch codeChallengeMethod {
	case "S256":
		codeChallenge = verifier.CodeChallengeS256()
	case "plain":
		codeChallenge = verifier.CodeChallengePlain()
	default:
		return "", "", "", errors.New("unknown Code Challenge Method: " + codeChallengeMethod)
	}

	codeVerifier = verifier.String()

	return codeChallengeMethod, codeVerifier, codeChallenge, nil
}

func extractStringSlice(data interface{}) ([]string, error) {
	if data == nil {
		return []string{}, nil
	}

	dataSlice, ok := data.([]interface{})
	if !ok {
		return []string{}, errors.New("not a slice")
	}

	dataStringSlice := []string{}

	for idx, dataElement := range dataSlice {
		dataString, ok := dataElement.(string)
		if !ok {
			return []string{}, errors.New("not a string element at index " + strconv.Itoa(idx))
		}

		dataStringSlice = append(dataStringSlice, dataString)
	}

	return dataStringSlice, nil
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
