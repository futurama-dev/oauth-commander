package client

import (
	"github.com/stretchr/testify/assert"
	"github.com/zalando/go-keyring"
	"testing"
)

func Test_load_missing(t *testing.T) {
	assert.DirExists(t, "../testdata/server")
	assert.NoDirExists(t, "../testdata/server/example_com")

	clients := load("../testdata/server/example_com")

	assert.Empty(t, clients)
}

func Test_load_empty(t *testing.T) {
	assert.DirExists(t, "../testdata/server")
	assert.DirExists(t, "../testdata/server/example_net")

	clients := load("../testdata/server/example_net")

	assert.Empty(t, clients)
}

func Test_load(t *testing.T) {
	assert.DirExists(t, "../testdata/server")
	assert.DirExists(t, "../testdata/server/example_org")

	clients := load("../testdata/server/example_org")

	assert.Len(t, clients, 2)
}

func TestClients_FindBySlug(t *testing.T) {
	clients := load("../testdata/server/example_org")

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

func TestClients_FindById(t *testing.T) {
	clients := load("../testdata/server/example_org")

	tests := []struct {
		name string
		id   string
		want bool
	}{
		{"found", "client_id_1", true},
		{"not found", "client_id_9", false},
		{"slug", "client_1", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ok := clients.FindById(tt.id)

			assert.Equalf(t, tt.want, ok, "FindBySlug(%v)", tt.id)

			if ok {
				assert.Equalf(t, tt.id, client.Id, "FindBySlug(%v)", tt.id)
			}
		})
	}
}

func TestClient_Secret(t *testing.T) {
	keyring.MockInit()

	client := Client{
		Slug:         "client_1",
		ServerSlug:   "example_org",
		Id:           "client_id_1",
		SecretHandle: "example_org_client_1",
	}

	secret, err := client.Secret()

	assert.Error(t, err)
	assert.Equal(t, keyring.ErrNotFound, err)

	err = client.SetSecret("sesame")

	assert.NoError(t, err)

	secret, err = client.Secret()

	assert.NoError(t, err)
	assert.Equal(t, "sesame", secret)
}
