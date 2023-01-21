package authorization

import (
	"errors"
	"fmt"
	"github.com/futurama-dev/oauth-commander/client"
	"github.com/futurama-dev/oauth-commander/oauth2"
	"github.com/futurama-dev/oauth-commander/session"
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
		log.Fatalln("not implemented")
	} else if req.Method == "GET" {
		authResponse, err := ProcessResponseUrl(responseUrl)

		if err != nil {
			errorResp, ok := err.(oauth2.ErrorResponse)
			if ok {
				errorResp.Println()
			} else {
				log.Fatalln(err)
			}
		} else {
			err = client.ProcessResponse(authResponse)
			if err != nil {
				log.Fatalln(err)
			}
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

func ProcessResponseUrl(responseUrl *url.URL) (oauth2.AuthorizationResponse, error) {
	query := responseUrl.Query()
	if !query.Has("state") {
		return oauth2.AuthorizationResponse{}, errors.New("state parameter missing")
	}

	state := query.Get("state")

	sessions, err := session.GetAuthorizationSessions()
	if err != nil {
		return oauth2.AuthorizationResponse{}, err
	}

	sess, found := sessions.FindByState(state)
	if !found {
		return oauth2.AuthorizationResponse{}, errors.New("no session found for state: " + state)
	}

	sess.CompletedAt = time.Now().Truncate(time.Second)
	err = session.SetAuthorizationSessions(sessions)
	if err != nil {
		return oauth2.AuthorizationResponse{}, err
	}

	if query.Has("error") {
		var errorUri *url.URL = nil

		if query.Has("error_uri") {
			errorUriStr := query.Get("error_uri")

			errorUri, err = url.Parse(errorUriStr)
			if err != nil {
				return oauth2.AuthorizationResponse{}, errors.New("invalid error_uri: " + errorUriStr)
			}
		}

		return oauth2.AuthorizationResponse{}, oauth2.ErrorResponse{
			ErrorCode:        query.Get("error"),
			ErrorDescription: query.Get("error_description"),
			ErrorUri:         errorUri,
			State:            query.Get("state"),
		}
	}

	return oauth2.AuthorizationResponse{
		State:   query.Get("state"),
		Code:    query.Get("code"),
		Token:   query.Get("token"),
		IdToken: query.Get("id_token"),
	}, nil
}
