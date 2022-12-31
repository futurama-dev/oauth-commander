package oidc

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const WellKnown = "/.well-known/openid-configuration"

func ValidateIssuer(issuer string) (*url.URL, error) {
	if len(issuer) == 0 {
		return nil, errors.New("issuer cannot be empty")
	}

	issuerUrl, err := url.Parse(issuer)

	if err != nil {
		return issuerUrl, err
	}

	if issuerUrl.Scheme != "https" {
		return issuerUrl, errors.New("issuer scheme must be https")
	}

	if len(issuerUrl.RawQuery) != 0 || issuerUrl.ForceQuery {
		return issuerUrl, errors.New("issuer cannot have a query")
	}

	if len(issuerUrl.Fragment) != 0 || strings.HasSuffix(issuer, "#") {
		return issuerUrl, errors.New("issuer cannot have a fragment")
	}
	return issuerUrl, nil
}

func BuildDiscoveryUrl(issuer url.URL) url.URL {
	if strings.HasSuffix(issuer.Path, "/") {
		issuer.Path += WellKnown[1:]
	} else {
		issuer.Path += WellKnown
	}

	return issuer
}

func FetchDiscovery(discoveryUrl string) (string, error) {
	resp, err := http.Get(discoveryUrl)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("failed with status: " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
