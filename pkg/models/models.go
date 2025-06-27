package models

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Abhishek-Omniful/OMS/mycontext"
	"github.com/Abhishek-Omniful/OMS/pkg/appinit"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/sqs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var logger = log.DefaultLogger()

type StoreCSV struct {
	FilePath string `json:"filePath"`
}

type BulkOrderRequest struct {
	FilePath string `json:"filePath"`
}

type Webhook struct {
	URL      string `json:"url" bson:"url"`
	TenantID int64  `json:"tenant_id" bson:"tenant_id"`
}

var mongoClient *mongo.Client
var err error
var client *awsS3.Client
var ctx context.Context
var ordersCollection *mongo.Collection
var webhookCollection *mongo.Collection
var publisher *sqs.Publisher

func init() {
	ctx = mycontext.GetContext()

	appinit.ConnectDB()
	ordersCollection = appinit.GetOrdersCollection()
	webhookCollection = appinit.GetWebhookCollection()

	appinit.ConnectS3()
	client = appinit.GetS3Client()

	appinit.SQSInit()
	newQueue := appinit.GetSqs()

	appinit.PublisherInit(newQueue)
	publisher = appinit.GetPublisher()

	appinit.InitConsumer()
	appinit.StartConsumer(ctx)

	go appinit.InitKafkaConsumer()
	time.Sleep(3 * time.Second)
	appinit.InitKafkaProducer()
	appinit.ConnectRedis()
	appinit.OrderRetryWorker()
}

func StoreInS3(s *StoreCSV) error {
	filepath := s.FilePath
	fileBytes := appinit.GetLocalCSV(filepath)

	bucketName := config.GetString(ctx, "s3.bucketName")
	filename := config.GetString(ctx, "s3.fileName")

	input := &awsS3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &filename,
		Body:   bytes.NewReader(fileBytes),
	}

	_, err := client.PutObject(ctx, input)
	if err != nil {
		logger.Error(i18n.Translate(ctx, err.Error()))
		return errors.New(i18n.Translate(ctx, "failed to upload to s3"))
	}
	logger.Infof(i18n.Translate(ctx, "File uploaded to S3!"))
	return nil
}

func ValidateS3Path_PushToSQS(req *BulkOrderRequest) error {
	logger.Infof(i18n.Translate(ctx, "Validating S3 path:"))
	filePath := req.FilePath

	if !strings.HasPrefix(filePath, "s3://") {
		return errors.New(i18n.Translate(ctx, "invalid S3 path format: must start with s3://"))
	}

	path := strings.TrimPrefix(filePath, "s3://")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 {
		return errors.New(i18n.Translate(ctx, "invalid S3 path: must be in s3://bucket/key format"))
	}

	bucket := parts[0]
	key := parts[1]

	logger.Infof(bucket, key)

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
	payload := fmt.Sprintf(`{"bucket":"%s", "key":"%s"}`, bucket, key)
	msg := &sqs.Message{
		Value: []byte(payload),
	}

	err = publisher.Publish(ctx, msg)
	if err != nil {
		logger.Error(i18n.Translate(ctx, "Failed to publish message to SQS:"), i18n.Translate(ctx, err.Error()))
		return err
	}
	logger.Infof(i18n.Translate(ctx, "Message successfully published to SQS"))
	return nil
}

func CreateWebhook(req *Webhook) error {
	if req.URL == "" || req.TenantID <= 0 {
		return errors.New(i18n.Translate(ctx, "invalid webhook request"))
	}

	webhook := &Webhook{
		URL:      req.URL,
		TenantID: req.TenantID,
	}

	_, err := webhookCollection.InsertOne(ctx, webhook)
	if err != nil {
		logger.Error(i18n.Translate(ctx, "Failed to create webhook:"), i18n.Translate(ctx, err.Error()))
		return err
	}

	logger.Infof(i18n.Translate(ctx, "Webhook created successfully!"))
	return nil
}

func ListWebhooks() ([]Webhook, error) {
	cursor, err := webhookCollection.Find(ctx, bson.M{})
	if err != nil {
		logger.Error(i18n.Translate(ctx, "Failed to list webhooks:"), i18n.Translate(ctx, err.Error()))
		return nil, err
	}
	defer cursor.Close(ctx)

	var webhooks []Webhook
	if err := cursor.All(ctx, &webhooks); err != nil {
		logger.Error(i18n.Translate(ctx, "Failed to decode webhooks:"), i18n.Translate(ctx, err.Error()))
		return nil, err
	}

	return webhooks, nil
}
