package appinit

import (
	"log"

	"github.com/Abhishek-Omniful/OMS/pkg/helper"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"

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

func GetLocalCSV(filepath string) []byte {
	fileBytes, _ := helper.GetLocalCSV(filepath)
	return fileBytes
}
