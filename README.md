AWS Resources:

Lambda:

ProcessFileFunction

S3:

itg-hack-day-file-uploads
itg-hack-day-file-uploads-processed

API Gateway:

itg-file-processor-staging

IAM:

Role: GoLambda
Policy: itg-hack-day-file-uploads-s3

SNS:

Topic: itg-hack-day-file-uploads
Subscription: itg-hack-day-file-upload-event
