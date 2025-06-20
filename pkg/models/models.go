// package models

// import (
// 	"bytes"
// 	"context"
// 	"errors"
// 	"log"
// 	"strings"

// 	"github.com/Abhishek-Omniful/OMS/mycontext"
// 	"github.com/Abhishek-Omniful/OMS/pkg/appinit"
// 	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
// 	"github.com/omniful/go_commons/config"
// 	"github.com/omniful/go_commons/sqs"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// type StoreCSV struct {
// 	FilePath string `json:"filePath"`
// }

// type BulkOrderRequest struct {
// 	SellerID string `json:"sellerID"`
// 	FilePath string `json:"filePath"`
// }

// var mongoClinet *mongo.Client
// var err error
// var client *awsS3.Client //  this is being returned to me by s3.NewDefaultAWSS3Client() of gocommons
// var ctx context.Context
// var collection *mongo.Collection
// var publisher *sqs.Publisher

// func init() {

// 	ctx := mycontext.GetContext()
// 	dbname := config.GetString(ctx, "mongo.dbname")
// 	collectionName := config.GetString(ctx, "mongo.collectionName")

// 	// Initialize MongoDB client and collection
// 	collection, err = appinit.GetMongoCollection(dbname, collectionName)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Initialize S3 client
// 	client = appinit.GetS3Client()

// 	// appinit.GetSqs()                   // Initialize SQS client
// 	// publisher = appinit.GetPublisher() // Initialize Publisher
// }

// func StoreInS3(s *StoreCSV) error {
// 	filepath := s.FilePath
// 	fileBytes := appinit.GetLocalCSV(filepath)
// 	bucketName := config.GetString(ctx, "s3.bucketName")
// 	filename := config.GetString(ctx, "s3.fileName")

// 	input := &awsS3.PutObjectInput{
// 		Bucket: &bucketName,
// 		Key:    &filename,
// 		Body:   bytes.NewReader(fileBytes),
// 	}

// 	_, err := client.PutObject(ctx, input)
// 	if err != nil {
// 		return errors.New("failed to upload to s3")
// 	}
// 	log.Println("File uploaded to S3!")
// 	return nil
// }

// func ValidateS3Path(req *BulkOrderRequest) error {
// 	filePath := req.FilePath

// 	if !strings.HasPrefix(filePath, "s3://") {
// 		return errors.New("invalid S3 path format: must start with s3://")
// 	}

// 	path := strings.TrimPrefix(filePath, "s3://")
// 	parts := strings.SplitN(path, "/", 2)
// 	if len(parts) != 2 {
// 		return errors.New("invalid S3 path: must be in s3://bucket/key format")
// 	}

// 	bucket := parts[0]
// 	key := parts[1]

// 	_, err := client.HeadObject(ctx, &awsS3.HeadObjectInput{
// 		Bucket: &bucket,
// 		Key:    &key,
// 	})

// 	if err != nil {
// 		return errors.New("file does not exist at specified S3 path")
// 	}

// 	return nil
// }

package models

import (
	"bytes"
	"context"
	"errors"
	"log"
	"strings"

	"github.com/Abhishek-Omniful/OMS/mycontext"
	"github.com/Abhishek-Omniful/OMS/pkg/appinit"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/omniful/go_commons/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type StoreCSV struct {
	FilePath string `json:"filePath"`
}

type BulkOrderRequest struct {
	SellerID string `json:"sellerID"`
	FilePath string `json:"filePath"`
}

var mongoClinet *mongo.Client
var err error
var client *awsS3.Client //  this is being returned to me by s3.NewDefaultAWSS3Client() of gocommons
var ctx context.Context
var collection *mongo.Collection

func init() {
	ctx = mycontext.GetContext()
	dbname := config.GetString(ctx, "mongo.dbname")
	collectionName := config.GetString(ctx, "mongo.collectionName")
    
	// Initialize MongoDB client and collection
	_, err = appinit.GetDB()
	if err != nil {
		log.Fatal(err)
	}
	collection, err = appinit.GetMongoCollection(dbname, collectionName)
	if err != nil {
		log.Fatal(err)
	}
	// Initialize S3 client
	client = appinit.GetS3Client()

     appinit.GetSqs()                   // Initialize SQS client
	// publisher = appinit.GetPublisher() // Initialize Publisher
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
		return errors.New("failed to upload to s3")
	}
	log.Println("File uploaded to S3!")
	return nil
}

func ValidateS3Path(req *BulkOrderRequest) error {
	filePath := req.FilePath

	if !strings.HasPrefix(filePath, "s3://") {
		return errors.New("invalid S3 path format: must start with s3://")
	}

	path := strings.TrimPrefix(filePath, "s3://")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 {
		return errors.New("invalid S3 path: must be in s3://bucket/key format")
	}

	bucket := parts[0]
	key := parts[1]

	_, err := client.HeadObject(ctx, &awsS3.HeadObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	//log.Println(location)

	if err != nil {
		return errors.New("file does not exist at specified S3 path")
	}

	return nil
}
