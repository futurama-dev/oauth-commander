package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIssuerToName(t *testing.T) {
	tests := []struct {
		name      string
		issuer    string
		want      string
		wantError bool
	}{
		{"empty", "", "", true},
		{"invalid", "abc_123", "abc_123", true},
		{"host", "https://example.com", "example_com", false},
		{"host and empty path", "https://example.com/", "example_com", false},
		{"host and empty paths", "https://example.com///", "example_com", false},
		{"host and port", "https://example.com:8080", "example_com_8080", false},
		{"host, port and empty path", "https://example.com:8080/", "example_com_8080", false},
		{"host, port and path", "https://example.com:8080/some/path", "example_com_8080_some_path", false},
		{"no schema", "example.com:8080/some/path", "example_com_8080_some_path", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IssuerToSlug(tt.issuer)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				if got != tt.want {
					t.Errorf("IssuerToSlug() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_load(t *testing.T) {
	assert.DirExists(t, "../testdata/server")

	servers := load("../testdata/server")

	assert.Len(t, servers, 2)
}

func Test_load_empty(t *testing.T) {
	assert.DirExists(t, "../testdata/server_empty")

	servers := load("../testdata/server_empty")

	assert.Empty(t, servers)
}

func TestServers_FindBySlug(t *testing.T) {
	assert.DirExists(t, "../testdata/server")

	servers := load("../testdata/server")

	tests := []struct {
		name string
		slug string
		want bool
	}{
		{"found", "example_com", true},
		{"not found", "example_org", false},
		{"issuer", "https://example.com/", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, ok := servers.FindBySlug(tt.slug)
			assert.Equalf(t, tt.want, ok, "FindBySlug(%v)", tt.slug)

			if ok {
				assert.Equalf(t, tt.slug, server.Slug, "FindBySlug(%v)", tt.slug)
			}
		})
	}
}

func TestServers_FindByIssuer(t *testing.T) {
	assert.DirExists(t, "../testdata/server")

	servers := load("../testdata/server")

	tests := []struct {
		name   string
		issuer string
		want   bool
	}{
		{"found", "https://example.com/", true},
		{"not found", "https://example.org/", false},
		{"slug", "example_com", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, ok := servers.FindByIssuer(tt.issuer)
			assert.Equalf(t, tt.want, ok, "FindBySlug(%v)", tt.issuer)

			if ok {
				assert.Equalf(t, tt.issuer, server.Metadata["issuer"], "FindBySlug(%v)", tt.issuer)
			}
		})
	}
}
