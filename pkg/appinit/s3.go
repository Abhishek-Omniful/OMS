package appinit

import (
	"bytes"
	"log"

	"github.com/Abhishek-Omniful/OMS/mycontext"
	"github.com/Abhishek-Omniful/OMS/pkg/helper"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/s3"
)

func connectS3() (*awsS3.Client, error) {
	//s3client creation
	log.Println("Connecting to s3")
	s3Client, err := s3.NewDefaultAWSS3Client()
	if err != nil {
		return nil, err
	}
	log.Println("Successfully Connected to s3")
	return s3Client, nil
}

func GetS3Client() *awsS3.Client {
	client, _ := connectS3()
	return client
}

func getLocalCSV() []byte {
	fileBytes, _ := helper.GetLocalCSV()
	return fileBytes
}

func UploadToLocalStack() {
	client := GetS3Client()
	ctx := mycontext.GetContext()
	fileBytes := getLocalCSV()

	bucketName := config.GetString(ctx, "s3.bucketName")
	filename := config.GetString(ctx, "s3.fileName")

	input := &awsS3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &filename,
		Body:   bytes.NewReader(fileBytes),
	}

	_, err := client.PutObject(ctx, input)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(" File uploaded to S3!")
}
