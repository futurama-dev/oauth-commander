package oauth2

import "github.com/futurama-dev/oauth-commander/discovery"

const WellKnown = "/.well-known/oauth-authorization-server"

func BuildDiscoveryUrl(issuer string) (string, error) {
	return discovery.BuildDiscoveryUrl(issuer, WellKnown)
}
