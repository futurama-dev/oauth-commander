package client

import (
	"github.com/stretchr/testify/assert"
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
