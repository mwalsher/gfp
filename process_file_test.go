package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestImageHandler(t *testing.T) {
	tests := []struct {
		request events.APIGatewayProxyRequest
		expect  string
		err     error
	}{
		{
			// Test that the handler responds with the correct response
			// when a valid name is provided in the HTTP body
			request: events.APIGatewayProxyRequest{},
			expect:  "BAM!!!!",
			err:     nil,
		},
		// {
		// 	// Test that the handler responds ErrNameNotProvided
		// 	// when no name is provided in the HTTP body
		// 	request: events.APIGatewayProxyRequest{Body: ""},
		// 	expect:  "",
		// 	err:     main.ErrNameNotProvided,
		// },
	}

	for _, test := range tests {
		response, err := imageHandler(test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.Body)
	}
}
