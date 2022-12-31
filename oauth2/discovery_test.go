package oauth2

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestBuildDiscoveryUrl(t *testing.T) {
	tests := []struct {
		name   string
		issuer string
		want   string
	}{
		{"No path", "https://example.com", "https://example.com/.well-known/oauth-authorization-server"},
		{"Empty path", "https://example.com/", "https://example.com/.well-known/oauth-authorization-server"},
		{"Path no trailing slash", "https://example.com/some/path", "https://example.com/some/path/.well-known/oauth-authorization-server"},
		{"Path with trailing slash", "https://example.com/some/path/", "https://example.com/some/path/.well-known/oauth-authorization-server"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildDiscoveryUrl(tt.issuer)

			assert.NoError(t, err)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildDiscoveryUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
