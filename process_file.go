// Examples: https://github.com/aws/aws-lambda-go
package main

import (
	// "net/http"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	// "github.com/aws/aws-sdk-go/service/s3"
)

/*
TODO:

Go

https://docs.aws.amazon.com/apigateway/latest/developerguide/integrating-api-with-aws-services-s3.html

Direct S3 upload via signed url examples:

https://www.netlify.com/blog/2016/11/17/serverless-file-uploads/
https://sookocheff.com/post/api/uploading-large-payloads-through-api-gateway/

1. Get file upload working
2. Return file name
3. Upload to S3
4. Process images

Rails

1. Make an HTTP call to Lambda to print out the result (via jQuery)
2. Change Attachments JS to create

JS to send file to Lambda endpoint
3. On file upload success, create Attachment with file info

*/

const (
	InputBucket  = "itg-hack-day-file-uploads"
	OutputBucket = "itg-hack-day-file-uploads-processed"
)

type FileProcessService struct {
	s3Client *s3.S3
}

func (svc *FileProcessService) GetFileMetadata(bucket string, key string) (*s3.HeadObjectOutput, error) {
	input := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	result, err := svc.s3Client.HeadObject(input)
	if err != nil {
		return nil, err
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
				return nil, aerr
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
			return nil, err
		}
	}

	return result, nil
}

func (svc *FileProcessService) GetFileContents(bucket string, key string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	result, err := svc.s3Client.GetObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				fmt.Println(s3.ErrCodeNoSuchKey, aerr.Error())
				return nil, err
			default:
				fmt.Println(aerr.Error())
				return nil, err
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
			return nil, err
		}
	}

	defer result.Body.Close()
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return body, nil
}

func (svc *FileProcessService) WriteFile(bucket string, key string, content []byte) (*s3.PutObjectOutput, error) {
	input := &s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(bytes.NewReader(content)),
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	result, err := svc.s3Client.PutObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
				return nil, err
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
			return nil, err
		}
	}

	return result, nil
}

func NewFileProcessService() *FileProcessService {
	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	// Create S3 service client
	s3Client := s3.New(sess)

	return &FileProcessService{
		s3Client: s3Client,
	}
}

func processObject(
	svc *FileProcessService,
	s3Object *events.S3Object,
) (string, error) {
	content, err := svc.GetFileContents(InputBucket, s3Object.Key)
	_, err = svc.WriteFile(OutputBucket, s3Object.Key, content)
	return string(content), err
}

// TODO Pass context to e.g. GetObjectWithContext
func handler(ctx context.Context, s3Event events.S3Event) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	svc := NewFileProcessService()

	var result string
	var s3Object *events.S3Object
	var err error
	for _, record := range s3Event.Records {
		s3Entity := record.S3
		s3Object = &s3Entity.Object
		fmt.Printf("[%s - %s] Bucket = %s, Key = %s \n", record.EventSource, record.EventTime, s3Entity.Bucket.Name, s3Object.Key)
		result, err = processObject(svc, s3Object)
		if err != nil {
			return "", err
		}
	}

	return result, nil
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
