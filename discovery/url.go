package discovery

import (
	"strings"
)

func BuildDiscoveryUrl(issuer string, well_known string) (string, error) {
	issuerUrl, err := ValidateIssuer(issuer)

	if err != nil {
		return "", err
	}

	if strings.HasSuffix(issuerUrl.Path, "/") {
		issuerUrl.Path += well_known[1:]
	} else {
		issuerUrl.Path += well_known
	}

	return issuerUrl.String(), nil
}
