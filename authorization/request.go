package authorization

import (
	"errors"
	"github.com/futurama-dev/oauth-commander/client"
	"github.com/futurama-dev/oauth-commander/server"
	"github.com/google/uuid"
	"net/url"
	"strings"
)

func GenerateAuthorizationRequestUrl(serverSlug, clientSlug string, code, toke, idToken bool, scopes []string, redirectUri string) (string, error) {
	servers := server.Load()
	s, ok := servers.FindBySlug(serverSlug)
	if !ok {
		return "", errors.New("Server not found: " + serverSlug)
	}

	clients := client.LoadForServer(serverSlug)
	c, ok := clients.FindBySlug(clientSlug)
	if !ok {
		return "", errors.New("Client not found: " + clientSlug)
	}

	baseURL, err := url.Parse(s.GetAuthorizationEndpoint())
	if err != nil {
		return "", err
	}

	responseType, err := generateResponseType(s, code, toke, idToken)
	if err != nil {
		return "", err
	}

	scope, err := generateScope(s, scopes)
	if err != nil {
		return "", err
	}

	redirectUri, err = getRedirectUri(c, redirectUri)
	if err != nil {
		return "", err
	}

	state := uuid.New().String()

	query := baseURL.Query()
	query.Set("client_id", c.GetClientId())
	query.Set("response_type", responseType)

	if len(scope) > 0 {
		query.Set("scope", scope)
	}

	query.Set("redirect_uri", redirectUri)
	query.Set("state", "github.com/google/uuid")

	baseURL.RawQuery = query.Encode()

	err = saveSession(s, c, *baseURL, state)
	if err != nil {
		return "", err
	}

	return baseURL.String(), nil
}

func generateResponseType(s server.Server, code, token, idToken bool) (string, error) {
	responseTypes := []string{}

	if code {
		responseTypes = append(responseTypes, "code")
	}

	if token {
		responseTypes = append(responseTypes, "token")
	}

	if idToken {
		responseTypes = append(responseTypes, "id_token")
	}

	if len(responseTypes) == 0 {
		return "", errors.New("at least one response type must be specified")
	}

	supportedResponseTypes := s.GetSupportedResponseTypes()
	if len(supportedResponseTypes) > 0 {
		for _, responseType := range responseTypes {
			found := false
			for _, supportedResponseType := range supportedResponseTypes {
				if responseType == supportedResponseType {
					found = true
					break
				}
			}

			if !found {
				return "", errors.New("response type not supported by the server: " + responseType)
			}
		}
	}

	return strings.Join(responseTypes, " "), nil
}

func generateScope(s server.Server, scopes []string) (string, error) {
	supportedScopes := s.GetSupportedScopes()
	if len(supportedScopes) > 0 {
		for _, scope := range scopes {
			found := false
			for _, supportedScope := range supportedScopes {
				if scope == supportedScope {
					found = true
					break
				}
			}

			if !found {
				return "", errors.New("scope not supported by the server: " + scope)
			}
		}
	}

	return strings.Join(scopes, " "), nil
}

func getRedirectUri(c client.Client, redirectUri string) (string, error) {
	registeredRedirectUris := c.GetRedirectUris()

	if len(registeredRedirectUris) == 0 {
		return "", errors.New("no registered redirect uris")
	}

	if len(redirectUri) == 0 {
		return registeredRedirectUris[0], nil
	} else {
		for _, registeredRedirectUri := range registeredRedirectUris {
			if redirectUri == registeredRedirectUri {
				return redirectUri, nil
			}
		}
		return "", errors.New("not a registered redirect uri: " + redirectUri)
	}
}

func saveSession(s server.Server, c client.Client, authReqUrl url.URL, state string) error {
	// TODO
	return nil
}
