package client

import (
	"encoding/json"
	"errors"
	"github.com/futurama-dev/oauth-commander/config"
	"github.com/futurama-dev/oauth-commander/oauth2"
	"github.com/futurama-dev/oauth-commander/server"
	"github.com/futurama-dev/oauth-commander/session"
	"github.com/zalando/go-keyring"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	Slug         string
	ServerSlug   string    `yaml:"server_slug"`
	CreatedAt    time.Time `yaml:"created_at"`
	Id           string
	SecretHandle string   `yaml:"secret_handle"`
	RedirectURIs []string `yaml:"redirect_uris"`
}

func (c Client) GetClientId() string {
	return c.Id
}

func (c Client) GetRedirectUris() []string {
	if c.RedirectURIs == nil {
		return []string{}
	} else {
		return c.RedirectURIs
	}
}

func (c Client) GetServer() (server.Server, error) {
	servers := server.Load()

	svr, found := servers.FindBySlug(c.ServerSlug)
	if !found {
		return server.Server{}, errors.New("server slug not found: " + c.ServerSlug)
	}

	return svr, nil
}

func (c Client) ExchangeCode(code, redirectUri, codeVerifier string) (oauth2.AccessTokenResponse, error) {
	svr, err := c.GetServer()
	if err != nil {
		return oauth2.AccessTokenResponse{}, err
	}

	data := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {code},
		"redirect_uri": {redirectUri},
		"client_id":    {c.Id},
	}

	if codeVerifier != "" {
		data["code_verifier"] = []string{codeVerifier}
	}

	clientSecret, err := c.GetSecret()
	if err != nil {
		return oauth2.AccessTokenResponse{}, err
	}

	authMethod := svr.GetTokenEndpointAuthMethod()
	useBasicAuth := false
	switch authMethod {
	case "client_secret_basic":
		useBasicAuth = true
	case "client_secret_post":
		data["client_secret"] = []string{clientSecret}
	default:
		return oauth2.AccessTokenResponse{}, errors.New("unsupported token endpoint auth method: " + authMethod)
	}

	req, err := http.NewRequest("POST", svr.GetTokenEndpoint(), strings.NewReader(data.Encode()))
	if err != nil {
		return oauth2.AccessTokenResponse{}, err
	}

	if useBasicAuth {
		req.SetBasicAuth(c.GetClientId(), clientSecret)
	}

	httpClient := http.Client{}

	resp, err := httpClient.Do(req)
	if err != nil {
		return oauth2.AccessTokenResponse{}, err
	}

	var tr oauth2.AccessTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tr)
	if err != nil {
		return oauth2.AccessTokenResponse{}, err
	}

	return tr, nil
}

func ProcessResponse(ar oauth2.AuthorizationResponse) error {
	ar.Println()

	if ar.State == "" {
		// TODO allow for responses without a state: either use selected server and client or take server and client
		// slugs as argument then retrieve all pending sessions for that server and client and if exactly one then
		// use that session
		return errors.New("cannot process an authorization response without a state")
	}

	sessions, err := session.GetAuthorizationSessions()
	if err != nil {
		return err
	}

	s, found := sessions.FindByState(ar.State)
	if !found {
		return errors.New("state not found: " + ar.State)
	}

	if s.IsExpired() {
		return errors.New("state expired: " + ar.State)
	}

	if ar.HasCode() {
		cl, found := LoadForServer(s.ServerSlug).FindBySlug(s.ClientSlug)
		if !found {
			return errors.New("client slug not found: " + s.ClientSlug + ", under server: " + s.ServerSlug)
		}

		redirectUri, err := s.GetRedirectUri()
		if err != nil {
			return err
		}

		tr, err := cl.ExchangeCode(ar.Code, redirectUri, s.CodeVerifier)
		if err != nil {
			return err
		}

		tr.Println()
		// TODO save tokens
	}

	if ar.HasTokens() {
		// TODO save tokens
	}

	return nil
}

type Clients []Client

func Load() Clients {
	serverSlug := config.GetSelectedServer()

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
func (c Client) GetSecret() (string, error) {
	return keyring.Get("oauth-commander", c.SecretHandle)
}

func (c Client) SetSecret(secret string) error {
	return keyring.Set("oauth-commander", c.SecretHandle, secret)
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
