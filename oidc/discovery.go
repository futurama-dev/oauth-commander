package oidc

import (
	"github.com/futurama-dev/oauth-commander/discovery"
)

const WellKnown = "/.well-known/openid-configuration"

func BuildDiscoveryUrl(issuer string) (string, error) {
	return discovery.BuildDiscoveryUrl(issuer, WellKnown)
}
