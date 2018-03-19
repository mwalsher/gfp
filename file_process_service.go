package main

import (
	"context"
	"fmt"
	"image/jpeg"
	"io"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/nfnt/resize"
)

type FileProcessService struct {
	s3client     *s3.S3
	s3uploader   *s3manager.Uploader
	s3downloader *s3manager.Downloader
}

func NewFileProcessService() *FileProcessService {
	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	// Create S3 service clients
	s3client := s3.New(sess)
	s3downloader := s3manager.NewDownloaderWithClient(s3client)
	s3uploader := s3manager.NewUploaderWithClient(s3client)

	return &FileProcessService{
		s3client:     s3client,
		s3downloader: s3downloader,
		s3uploader:   s3uploader,
	}
}

func (svc *FileProcessService) GetFileMetadata(
	ctx context.Context,
	bucket string,
	key string,
) (*s3.HeadObjectOutput, error) {
	input := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	result, err := svc.s3client.HeadObjectWithContext(ctx, input)
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

func (svc *FileProcessService) GetFileContents(
	ctx context.Context,
	bucket string,
	key string,
) ([]byte, error) {
	r, err := svc.GetFile(ctx, bucket, key)
	if err != nil {
		return nil, err
	}

	defer r.Close()
	content, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return content, nil
}

func (svc *FileProcessService) GetFile(
	ctx context.Context,
	bucket string,
	key string,
) (io.ReadCloser, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	result, err := svc.s3client.GetObjectWithContext(ctx, input)
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

	return result.Body, nil
}

func (svc *FileProcessService) Resize(r io.ReadCloser) (io.Reader, error) {
	defer r.Close()

	// decode jpeg into image.Image
	img, err := jpeg.Decode(r)
	if err != nil {
		return nil, err
	}

	// resize to width 1000 using NearestNeighbor resampling and preserve aspect ratio
	imgr := resize.Resize(1000, 0, img, resize.NearestNeighbor)

	// Set up a pipe to write data directly into the Reader (unbuffered)
	pr, pw := io.Pipe()

	go func() {
		jpeg.Encode(pw, imgr, nil)
		pw.Close()
	}()

	return pr, nil

	// Alternative, buffered solution (memory inefficient)
	// create a buffer
	// var b bytes.Buffer
	// log("jpeg.Encode")
	// jpeg.Encode(&b, imgr, nil)
	// log("jpeg.Encode Done")
	// return &b, nil
}

func (svc *FileProcessService) WriteFile(
	ctx context.Context,
	bucket string,
	key string,
	r io.Reader,
) error {
	input := &s3manager.UploadInput{
		Body:   r,
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	_, err := svc.s3uploader.UploadWithContext(ctx, input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
				return err
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
			return err
		}
	}

	return nil
}

func (svc *FileProcessService) DeleteFile(ctx context.Context, bucket string, key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	_, err := svc.s3client.DeleteObjectWithContext(ctx, input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
				return err
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
			return err
		}
	}

	return nil
}
