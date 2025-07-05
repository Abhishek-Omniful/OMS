package models

import (
	"errors"
	"fmt"
	"strings"

	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/sqs"
)

var logger = log.DefaultLogger()

// func ValidateS3Path_PushToSQS(req *BulkOrderRequest) error {
var ValidateS3Path_PushToSQS = func(req *BulkOrderRequest) error {
	logger.Infof(i18n.Translate(ctx, "Validating S3 path:"))
	filePath := req.FilePath

	//sample filePath: "s3://my-bucket/my-file.csv"

	if !strings.HasPrefix(filePath, "s3://") {
		return errors.New(i18n.Translate(ctx, "invalid S3 path format: must start with s3://"))
	}

	path := strings.TrimPrefix(filePath, "s3://")
	parts := strings.SplitN(path, "/", 2) // split string on basis of "/" into atmost 2 parts
	if len(parts) != 2 {
		return errors.New(i18n.Translate(ctx, "invalid S3 path: must be in s3://bucket/key format"))
	}

	bucket := parts[0] //bucket name
	key := parts[1]    //filename

	logger.Infof(bucket, key)

	// check if the file exists in the specified S3 bucket by sending a Head request
	// if file exists, it will return metadata about the object, otherwise it will return an error
	_, err := client.HeadObject(ctx, &awsS3.HeadObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})

	if err != nil {
		logger.Error(i18n.Translate(ctx, err.Error()))
		return errors.New(i18n.Translate(ctx, "file does not exist at specified S3 path"))
	}
	logger.Infof(i18n.Translate(ctx, "S3 path is valid Successfully!"))
	logger.Infof(i18n.Translate(ctx, "Pushing to SQS..."))

	err = PushToSQS(bucket, key)
	if err != nil {
		logger.Error(i18n.Translate(ctx, "Failed to push to SQS:"), i18n.Translate(ctx, err.Error()))
	}
	logger.Infof(i18n.Translate(ctx, "Successfully pushed to SQS!"))
	return nil
}

func PushToSQS(bucket string, key string) error {

	//  {"bucket":"your-bucket-name", "key":"your-key-name"}
	payload := fmt.Sprintf(`{"bucket":"%s", "key":"%s"}`, bucket, key)

	msg := &sqs.Message{
		Value: []byte(payload),
	}

	err = sqsPublisher.Publish(ctx, msg)
	if err != nil {
		logger.Error(i18n.Translate(ctx, "Failed to publish message to SQS:"), i18n.Translate(ctx, err.Error()))
		return err
	}
	logger.Infof(i18n.Translate(ctx, "Message successfully published to SQS"))
	//after this the control will go to the consumer which will process the message
	return nil
}
