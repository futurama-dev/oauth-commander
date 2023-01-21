package oauth2

import (
	"fmt"
	"net/url"
)

type AuthorizationResponse struct {
	State   string
	Code    string
	Token   string
	IdToken string
}

func (ar AuthorizationResponse) HasCode() bool {
	return len(ar.Code) > 0
}

func (ar AuthorizationResponse) HasTokens() bool {
	return len(ar.Token) > 0 || len(ar.IdToken) > 0
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
	fmt.Println("error:", er.ErrorCode)
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

type AccessTokenResponse map[string]interface{}

func (atr AccessTokenResponse) Println() {
	for k, v := range atr {
		fmt.Println(k, v)
	}
}
