package client

import (
	"github.com/stretchr/testify/assert"
	"github.com/zalando/go-keyring"
	"testing"
)

func Test_load_missing(t *testing.T) {
	assert.DirExists(t, "../testdata/server")
	assert.NoDirExists(t, "../testdata/server/example_com")

	clients := load("example_com", "../testdata/server")

	assert.Empty(t, clients)
}

func Test_load_empty(t *testing.T) {
	assert.DirExists(t, "../testdata/server")
	assert.DirExists(t, "../testdata/server/example_net")

	clients := load("example_net", "../testdata/server")

	assert.Empty(t, clients)
}

func Test_load(t *testing.T) {
	assert.DirExists(t, "../testdata/server")
	assert.DirExists(t, "../testdata/server/example_org")

	clients := load("example_org", "../testdata/server")

	assert.Len(t, clients, 2)
}

func TestClients_FindBySlug(t *testing.T) {
	clients := Load()

	tests := []struct {
		name string
		slug string
		want bool
	}{
		{"found", "client_1", true},
		{"not found", "client_9", false},
		{"id", "client_id_1", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ok := clients.FindBySlug(tt.slug)

			assert.Equalf(t, tt.want, ok, "FindBySlug(%v)", tt.slug)

			if ok {
				assert.Equalf(t, tt.slug, client.Slug, "FindBySlug(%v)", tt.slug)
			}
		})
	}
}

func TestClient_Secret(t *testing.T) {
	keyring.MockInit()

	client := Client{
		Slug:         "client_1",
		Type:         "oidc",
		Id:           "client_id_1",
		SecretHandle: "example_org_client_1",
	}

	err := client.SetSecret("sesame")

	assert.NoError(t, err)

	secret, err := client.Secret()

	assert.NoError(t, err)
	assert.Equal(t, "sesame", secret)
}
