package main

import (
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	os.Setenv("API_KEY", "MockingAPIKey")
	tests := []struct {
		request events.APIGatewayProxyRequest
		err     error
	}{
		{
			// Test that the handler responds with the correct body
			request: events.APIGatewayProxyRequest{Body: `{"id": 18}`},
			err:     nil,
		},
		{
			// Test that the handler responds with the empty body
			request: events.APIGatewayProxyRequest{Body: ``},
			err:     nil,
		},
	}
	for _, test := range tests {
		_, err := Handler(test.request)
		assert.IsType(t, test.err, err)
	}
}
