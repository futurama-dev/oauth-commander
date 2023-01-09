package config

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestSession_IsExpired(t *testing.T) {
	SetDefaults()

	authReqParsedUrl, err := url.Parse("https://example.com/authorize?client_id=some_other_client_id&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fcallback&response_type=code&scope=email&state=6cac71ed-8051-47e0-b3ed-50543204f8d2")
	assert.NoError(t, err)

	s := NewSession(
		"6cac71ed-8051-47e0-b3ed-50543204f8d2",
		*authReqParsedUrl,
		"some_server",
		"some_client")
	assert.False(t, s.IsExpired())
}
