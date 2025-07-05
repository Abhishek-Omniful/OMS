package sqsService

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/Abhishek-Omniful/OMS/mycontext"
	csvProcessorService "github.com/Abhishek-Omniful/OMS/pkg/helper/csvProcessor"
	dbService "github.com/Abhishek-Omniful/OMS/pkg/integrations/db"
	s3 "github.com/Abhishek-Omniful/OMS/pkg/integrations/s3" // go commons s3 import
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/sqs"
	"go.mongodb.org/mongo-driver/mongo"
)

var consumer *sqs.Consumer
var logger = log.DefaultLogger()
var s3Client *awsS3.Client
var ordersCollection *mongo.Collection

func InitSqsConsumer() {
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
		sqsQueue,             //SQS queue to consume messages from
		numberOfWorker,       // Number of workers to process messages
		concurrencyPerWorker, //Number of messages each worker can process concurrently
		&queueHandler{},      //Struct implementing the logic to handle each message
		maxMessagesCount,     // Maximum number of messages to fetch in one batch
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
	s3Client = s3.GetS3Client()
	ordersCollection = dbService.GetOrdersCollection()
}

type queueHandler struct{}

func (h *queueHandler) Process(ctx context.Context, msgs *[]sqs.Message) (err error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf(i18n.Translate(ctx, "Panic in queueHandler.Process: %v\nStack trace:\n%s"), r, debug.Stack())
			err = fmt.Errorf("panic occurred while processing SQS messages: %v", r)
		}
	}()

	for _, msg := range *msgs {
		var payload struct {
			Bucket string `json:"bucket"` // S3 bucket name
			Key    string `json:"key"`    // file name in S3
		}

		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			logger.Errorf(i18n.Translate(ctx, "Invalid message payload: %v"), err)
			continue
		}

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

		err = csvProcessorService.ParseCSV(tmpFile, ctx, logger, ordersCollection)
		if err != nil {
			logger.Errorf(i18n.Translate(ctx, "Failed to parse CSV file: %v"), err)
			continue
		}
		logger.Infof(i18n.Translate(ctx, "CSV file parsed successfully: %s"), tmpFile)
	}
	return nil
}
