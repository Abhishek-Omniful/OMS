package appinit

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/Abhishek-Omniful/OMS/mycontext"
	// parse_csv "github.com/Abhishek-Omniful/OMS/pkg/helper"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/sqs"
)

var consumer *sqs.Consumer

func InitConsumer() {
	sqsQueue := GetSqs()
	ctx := mycontext.GetContext()

	numberOfWorker := config.GetUint64(ctx, "consumer.numberOfWorker")
	concurrencyPerWorker := config.GetUint64(ctx, "consumer.concurrencyPerWorker")
	maxMessagesCount := config.GetInt64(ctx, "consumer.maxMessagesCount")
	visibilityTimeout := config.GetInt64(ctx, "consumer.visibilityTimeout")
	isAsync := config.GetBool(ctx, "consumer.isAsync")
	sendBatchMessage := config.GetBool(ctx, "consumer.sendBatchMessage")

	consumer, err = sqs.NewConsumer(
		sqsQueue,
		numberOfWorker,
		concurrencyPerWorker,
		&queueHandler{}, // defined below
		maxMessagesCount,
		visibilityTimeout,
		isAsync,
		sendBatchMessage,
	)

	if err != nil {
		logger.Panicf("Failed to start SQS consumer: %v", err)
	}

	logger.Infof("SQS consumer initialized")
}

func StartConsumer(ctx context.Context) {
	consumer.Start(ctx)
	logger.Infof("SQS consumer started")
}

type queueHandler struct{}

func (h *queueHandler) Process(ctx context.Context, msgs *[]sqs.Message) error {
	//s3Client := GetS3Client()  we dont need to get it as it is initialized in same package

	if err != nil {
		logger.Errorf("Failed to create S3 client: %v", err)
		return err
	}
	for _, msg := range *msgs {
		// Parse message payload
		var payload struct {
			Bucket string `json:"bucket"`
			Key    string `json:"key"`
		}
		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			logger.Errorf("Invalid message payload: %v", err)
			continue
		}

		// Download from S3
		_, err := s3Client.GetObject(ctx, &awsS3.GetObjectInput{
			Bucket: aws.String(payload.Bucket),
			Key:    aws.String(payload.Key),
		})
		if err != nil {
			logger.Errorf("Failed to download S3 object: %v", err)
			continue
		}

		// Step 1: Download CSV to local temp file
		tmpFile := filepath.Join(os.TempDir(), filepath.Base(payload.Key))
		getObjOutput, err := s3Client.GetObject(ctx, &awsS3.GetObjectInput{
			Bucket: aws.String(payload.Bucket),
			Key:    aws.String(payload.Key),
		})
		if err != nil {
			logger.Errorf("failed to download CSV from S3: %v", err)
			continue
		}
		logger.Infof("Downloaded CSV from S3")

		defer getObjOutput.Body.Close()

		outFile, err := os.Create(tmpFile)
		if err != nil {
			logger.Errorf("failed to create temp file: %v", err)
			continue
		}
		logger.Infof("Created temp file to store downloaded CSV")

		defer outFile.Close()

		_, err = io.Copy(outFile, getObjOutput.Body)
		if err != nil {
			logger.Errorf("failed to write CSV data to file: %v", err)
			continue
		}

		logger.Infof("CSV data written to temp file: %s", tmpFile)
		logger.Infof("Starting to parse CSV file: %s", tmpFile)

		// Parse the CSV file
		err = ParseCSV(tmpFile, ctx, logger, OrdersCollection)
		if err != nil {
			logger.Errorf("failed to parse CSV file: %v", err)
			continue
		}
		logger.Infof("CSV file parsed successfully : %s", tmpFile)
	}

	return nil
}
