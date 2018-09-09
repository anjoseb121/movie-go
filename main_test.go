package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	os.Setenv("API_KEY", "MockingAPIKey")
	movies, err := Handler(Request{
		ID: 18,
	})
	assert.IsType(t, nil, err)
	assert.NotEqual(t, 0, len(movies))
}
