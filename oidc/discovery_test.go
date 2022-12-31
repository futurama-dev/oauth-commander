package oidc

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
		{"No path", "https://example.com", "https://example.com/.well-known/openid-configuration"},
		{"Empty path", "https://example.com/", "https://example.com/.well-known/openid-configuration"},
		{"Path no trailing slash", "https://example.com/some/path", "https://example.com/some/path/.well-known/openid-configuration"},
		{"Path with trailing slash", "https://example.com/some/path/", "https://example.com/some/path/.well-known/openid-configuration"},
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
