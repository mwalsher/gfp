package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

/*
TODO:

Go

1. Get file upload working
2.

Rails

1. Make an HTTP call to Lambda to print out the result (via jQuery)
2. Change Attachments JS to send file to Lambda endpoint
3. On file upload success, create Attachment with file info

*/

func imageHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "BAM!!!!",
	}, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(imageHandler)
}
