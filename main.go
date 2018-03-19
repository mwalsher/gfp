// Examples: https://github.com/aws/aws-lambda-go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	InputBucket  = "itg-hack-day-file-uploads"
	OutputBucket = "itg-hack-day-file-uploads-processed"
)

func processObject(
	ctx context.Context,
	svc *FileProcessService,
	key string,
) error {
	fr, err := svc.GetFile(ctx, InputBucket, key)
	if err != nil {
		return err
	}
	defer svc.DeleteFile(ctx, InputBucket, key)
	r, err := svc.Resize(fr)
	if err != nil {
		return err
	}
	err = svc.WriteFile(ctx, OutputBucket, key, r)
	return err
}

func log(m string) {
	fmt.Println(m)
}

func handler(ctx context.Context, s3Event events.S3Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	svc := NewFileProcessService()

	var s3Object *events.S3Object
	var err error
	for _, record := range s3Event.Records {
		s3Entity := record.S3
		s3Object = &s3Entity.Object
		fmt.Printf("[%s - %s] Bucket = %s, Key = %s \n", record.EventSource, record.EventTime, s3Entity.Bucket.Name, s3Object.Key)
		err = processObject(ctx, svc, s3Object.Key)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}

// func handler(ctx context.Context, snsEvent events.SNSEvent) {
// 	for _, record := range snsEvent.Records {
// 		snsRecord := record.SNS
//
// 		fmt.Printf("[%s %s] Message = %s \n", record.EventSource, snsRecord.Timestamp, snsRecord.Message)
// 	}
// }

// func fileHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
// 	return events.APIGatewayProxyResponse{
// 		StatusCode: http.StatusOK,
// 		Body:       "BAM!!!!",
// 	}, nil
// }
//
