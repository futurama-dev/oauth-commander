package discovery

import (
	"errors"
	"net/url"
	"strings"
)

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
