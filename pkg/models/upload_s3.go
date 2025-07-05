package models

import (
	"bytes"
	"errors"

	getlocalcsv "github.com/Abhishek-Omniful/OMS/pkg/helper/getLocalCSV"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/i18n"
)



func StoreInS3(s *StoreCSV) error {
	filepath := s.FilePath                         // file path user provided
	fileBytes := getlocalcsv.GetLocalCSV(filepath) //read file from given path and store in fileBytes

	bucketName := config.GetString(ctx, "s3.bucketName")
	filename := config.GetString(ctx, "s3.fileName")

	input := &awsS3.PutObjectInput{ // prepare input for s3 upload
		Bucket: &bucketName,
		Key:    &filename,
		Body:   bytes.NewReader(fileBytes),
	}

	_, err := client.PutObject(ctx, input) // uploads the input to s3
	if err != nil {
		logger.Error(i18n.Translate(ctx, err.Error()))
		return errors.New(i18n.Translate(ctx, "failed to upload to s3"))
	}
	logger.Infof(i18n.Translate(ctx, "File uploaded to S3!"))
	return nil
}
