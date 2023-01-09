package authorization

import (
	"errors"
	"fmt"
	"github.com/futurama-dev/oauth-commander/config"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

func Listen() error {
	http.HandleFunc("/callback", oauthCallback)
	return http.ListenAndServe(":"+strconv.Itoa(8765), nil)
}

func oauthCallback(resp http.ResponseWriter, req *http.Request) {
	responseUrl := req.URL
	fmt.Println("Received response:", responseUrl.String())
	fmt.Println("HTTP method:", req.Method)
	if req.Method == "POST" {
		fmt.Println("Content-Type:", req.Header.Get("Content-Type"))
		fmt.Println("Content-Length:", req.ContentLength)
		err := req.ParseForm()
		if err != nil {
			log.Fatalln(err)
		}
		for key, value := range req.Form {
			fmt.Println(key, "=", value)
		}
	} else if req.Method == "GET" {
		authResponse, err := ProcessResponseUrl(responseUrl)

		if err != nil {
			errorResp, ok := err.(ErrorResponse)
			if ok {
				errorResp.Println()
			} else {
				log.Fatalln(err)
			}
		} else {
			// TODO handle successful response
			// exchange code
			// save all tokens
			authResponse.Println()
		}
	} else {
		log.Fatalln("HTTP method not supported: ", req.Method)
	}

	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	resp.WriteHeader(http.StatusOK)
	_, err := resp.Write([]byte("Response received.\n\nYou can close this window now.\n"))
	if err != nil {
		log.Fatalln(err)
	}

	os.Exit(0)
}

type AuthorizationResponse struct {
	State   string
	Code    string
	Token   string
	IdToken string
}

func (ar AuthorizationResponse) Println() {
	if ar.State != "" {
		fmt.Println("state:", ar.State)
	}
	if ar.Code != "" {
		fmt.Println("code:", ar.Code[:5]+"...")
	}
	if ar.Token != "" {
		fmt.Println("token:", ar.Token[:5]+"...")
	}
	if ar.IdToken != "" {
		fmt.Println("id_token:", ar.IdToken)
	}
}

type ErrorResponse struct {
	ErrorCode        string
	ErrorDescription string
	ErrorUri         *url.URL
	State            string
}

func (er ErrorResponse) Error() string {
	return er.ErrorCode
}

func (er ErrorResponse) Println() {
	fmt.Println("errpr:", er.ErrorCode)
	if er.ErrorDescription != "" {
		fmt.Println("error_description", er.ErrorDescription)
	}
	if er.ErrorUri != nil {
		fmt.Println("error_uri", er.ErrorUri.String())
	}
	if er.State != "" {
		fmt.Println("state", er.State)
	}
}

func ProcessResponseUrl(responseUrl *url.URL) (AuthorizationResponse, error) {
	query := responseUrl.Query()
	if !query.Has("state") {
		return AuthorizationResponse{}, errors.New("state parameter missing")
	}

	state := query.Get("state")

	sessions, err := config.GetAuthorizationSessions()
	if err != nil {
		return AuthorizationResponse{}, err
	}

	sess, found := sessions.FindByState(state)
	if !found {
		return AuthorizationResponse{}, errors.New("no session found for state: " + state)
	}

	sess.CompletedAt = time.Now().Truncate(time.Second)
	err = config.SetAuthorizationSessions(sessions)
	if err != nil {
		return AuthorizationResponse{}, err
	}

	if query.Has("error") {
		var errorUri *url.URL = nil

		if query.Has("error_uri") {
			errorUriStr := query.Get("error_uri")

			errorUri, err = url.Parse(errorUriStr)
			if err != nil {
				return AuthorizationResponse{}, errors.New("invalid error_uri: " + errorUriStr)
			}
		}

		return AuthorizationResponse{}, ErrorResponse{
			ErrorCode:        query.Get("error"),
			ErrorDescription: query.Get("error_description"),
			ErrorUri:         errorUri,
			State:            query.Get("state"),
		}
	}

	return AuthorizationResponse{
		State:   query.Get("state"),
		Code:    query.Get("code"),
		Token:   query.Get("token"),
		IdToken: query.Get("id_token"),
	}, nil
}
