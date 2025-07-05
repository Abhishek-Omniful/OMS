package models

import (
	"context"

	"github.com/Abhishek-Omniful/OMS/mycontext"
	dbService "github.com/Abhishek-Omniful/OMS/pkg/integrations/db"
	s3Service "github.com/Abhishek-Omniful/OMS/pkg/integrations/s3"
	sqsService "github.com/Abhishek-Omniful/OMS/pkg/integrations/sqs"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/omniful/go_commons/sqs"
	"go.mongodb.org/mongo-driver/mongo"
)

type BulkOrderRequest struct {
	FilePath string `json:"filePath"`
}

type StoreCSV struct {
	FilePath string `json:"filePath"`
}

type Webhook struct {
	URL      string `json:"url" bson:"url"`
	TenantID int64  `json:"tenant_id" bson:"tenant_id"`
}

var err error
var client *awsS3.Client
var ctx context.Context
var webhookCollection *mongo.Collection
var sqsPublisher *sqs.Publisher

func init() {
	client = s3Service.GetS3Client()
	ctx = mycontext.GetContext()
	webhookCollection = dbService.GetWebhookCollection()
	sqsPublisher = sqsService.GetSQSPublisher()
}
