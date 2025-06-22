package appinit

import (
	"log"

	"github.com/Abhishek-Omniful/OMS/pkg/helper"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/omniful/go_commons/s3"
)

var s3Client *awsS3.Client

func ConnectS3() {
	//s3client creation

	log.Println("Connecting to s3")
	s3Client, err = s3.NewDefaultAWSS3Client()
	if err != nil {
		log.Fatal("Error connecting to s3:", err)
		return
	}
	log.Println("Successfully Connected to s3")
}

func GetS3Client() *awsS3.Client {
	return s3Client
}

func GetLocalCSV(filepath string) []byte {
	fileBytes, _ := helper.GetLocalCSV(filepath)
	return fileBytes
}
