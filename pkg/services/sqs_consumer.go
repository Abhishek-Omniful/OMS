package services

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/Abhishek-Omniful/OMS/mycontext"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/sqs"
)

var consumer *sqs.Consumer

func InitConsumer() {
	ctx := mycontext.GetContext()
	sqsQueue := GetSqs()

	numberOfWorker := config.GetUint64(ctx, "consumer.numberOfWorker")
	concurrencyPerWorker := config.GetUint64(ctx, "consumer.concurrencyPerWorker")
	maxMessagesCount := config.GetInt64(ctx, "consumer.maxMessagesCount")
	visibilityTimeout := config.GetInt64(ctx, "consumer.visibilityTimeout")
	isAsync := config.GetBool(ctx, "consumer.isAsync")
	sendBatchMessage := config.GetBool(ctx, "consumer.sendBatchMessage")

	var err error
	consumer, err = sqs.NewConsumer(
		sqsQueue,
		numberOfWorker,
		concurrencyPerWorker,
		&queueHandler{},
		maxMessagesCount,
		visibilityTimeout,
		isAsync,
		sendBatchMessage,
	)
	if err != nil {
		logger.Panicf(i18n.Translate(ctx, "Failed to start SQS consumer: %v"), err)
	}
	logger.Infof(i18n.Translate(ctx, "SQS consumer initialized"))
}

func StartConsumer(ctx context.Context) {
	consumer.Start(ctx)
	logger.Infof(i18n.Translate(ctx, "SQS consumer started"))
}

type queueHandler struct{}

func (h *queueHandler) Process(ctx context.Context, msgs *[]sqs.Message) error {
	for _, msg := range *msgs {
		var payload struct {
			Bucket string `json:"bucket"`
			Key    string `json:"key"`
		}

		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			logger.Errorf(i18n.Translate(ctx, "Invalid message payload: %v"), err)
			continue
		}

		// Step 1: Download object from S3
		tmpFile := filepath.Join(os.TempDir(), filepath.Base(payload.Key))
		getObjOutput, err := s3Client.GetObject(ctx, &awsS3.GetObjectInput{
			Bucket: aws.String(payload.Bucket),
			Key:    aws.String(payload.Key),
		})
		if err != nil {
			logger.Errorf(i18n.Translate(ctx, "Failed to download CSV from S3: %v"), err)
			continue
		}
		defer getObjOutput.Body.Close()

		logger.Infof(i18n.Translate(ctx, "Downloaded CSV from S3"))

		// Step 2: Write to local temp file
		outFile, err := os.Create(tmpFile)
		if err != nil {
			logger.Errorf(i18n.Translate(ctx, "Failed to create temp file: %v"), err)
			continue
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, getObjOutput.Body)
		if err != nil {
			logger.Errorf(i18n.Translate(ctx, "Failed to write CSV data to temp file: %v"), err)
			continue
		}

		logger.Infof(i18n.Translate(ctx, "CSV data written to temp file: %s"), tmpFile)
		logger.Infof(i18n.Translate(ctx, "Starting to parse CSV file: %s"), tmpFile)

		// Step 3: Parse CSV
		err = ParseCSV(tmpFile, ctx, logger, OrdersCollection)
		if err != nil {
			logger.Errorf(i18n.Translate(ctx, "Failed to parse CSV file: %v"), err)
			continue
		}
		logger.Infof(i18n.Translate(ctx, "CSV file parsed successfully: %s"), tmpFile)
	}
	return nil
}
