package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
	"time"
)

func Test_write(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "servers*")
	if err != nil {
		log.Fatalln(err)
	}

	defer os.RemoveAll(tempDir)

	testServer := Server{
		Slug:      "example_com",
		Type:      "oidc",
		CreatedAt: time.Now().Truncate(time.Second),
		Metadata:  map[string]any{},
	}

	// save tesServer
	err = write(testServer, false, tempDir)
	require.NoError(t, err)

	// load testServer
	results := load(tempDir)
	assert.Len(t, results, 1)

	assert.Equal(t, testServer, results[0])

	// save testServer onto existing testServer
	err = write(testServer, false, tempDir)
	require.Error(t, err)

	// update same testServer onto existing testServer
	err = write(testServer, true, tempDir)
	require.NoError(t, err)

	// update testServer onto existing testServer
	testServer.CreatedAt = testServer.CreatedAt.Add(time.Hour)
	err = write(testServer, true, tempDir)

	results = load(tempDir)
	assert.Len(t, results, 1)

	assert.Equal(t, testServer, results[0])

	testServer2 := Server{
		Slug:      "example_org",
		Type:      "oidc",
		CreatedAt: time.Now().Truncate(time.Second),
		Metadata:  map[string]any{},
	}

	err = write(testServer2, false, tempDir)
	require.NoError(t, err)

	results = load(tempDir)
	assert.Len(t, results, 2)
}
