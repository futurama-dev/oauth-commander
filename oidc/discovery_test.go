package oidc

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"reflect"
	"testing"
)

func TestValidateIssuer(t *testing.T) {
	tests := []struct {
		name    string
		issuer  string
		wantErr bool
	}{
		{"empty", "", true},
		{"invalid URL", "sdfgsdfgsdfgsdfg", true},
		{"HTTP", "http://example.com", true},
		{"Empty Query", "https://example.com/path?", true},
		{"Query", "https://example.com/path?q=1", true},
		{"Empty Fragment", "https://example.com/path#", true},
		{"Fragment", "https://example.com/path#frag", true},
		{"Good no path", "https://example.com", false},
		{"Good empty path", "https://example.com/", false},
		{"Good with path", "https://example.com/some/path", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ValidateIssuer(tt.issuer); (err != nil) != tt.wantErr {
				t.Errorf("ValidateIssuer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

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
			issuerUrl, err := url.Parse(tt.issuer)
			assert.NoError(t, err)
			got := BuildDiscoveryUrl(*issuerUrl)
			if !reflect.DeepEqual(got.String(), tt.want) {
				t.Errorf("BuildDiscoveryUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
