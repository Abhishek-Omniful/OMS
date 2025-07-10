package models

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateS3Path_PushToSQS(t *testing.T) {
	tests := []struct {
		name             string
		filePath         string
		mockValidateFunc func(string, string) bool
		mockPushFunc     func(string, string) error
		expectedErr      bool
	}{
		{
			name:     "Valid path and success",
			filePath: "s3://bucket/key.csv",
			mockValidateFunc: func(bucket, key string) bool {
				return true
			},
			mockPushFunc: func(bucket, key string) error {
				return nil
			},
			expectedErr: false,
		},
		{
			name:             "Invalid prefix",
			filePath:         "https://bucket/key.csv",
			expectedErr:      true,
			mockValidateFunc: nil,
			mockPushFunc:     nil,
		},
		{
			name:             "Invalid format",
			filePath:         "s3://bucketonly",
			expectedErr:      true,
			mockValidateFunc: nil,
			mockPushFunc:     nil,
		},
		{
			name:     "S3 file not found",
			filePath: "s3://bucket/key.csv",
			mockValidateFunc: func(bucket, key string) bool {
				return false
			},
			mockPushFunc: func(bucket, key string) error {
				return nil // won't be reached
			},
			expectedErr: true,
		},
		{
			name:     "SQS push fails",
			filePath: "s3://bucket/key.csv",
			mockValidateFunc: func(bucket, key string) bool {
				return true
			},
			mockPushFunc: func(bucket, key string) error {
				return errors.New("sqs error")
			},
			expectedErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx = context.Background()
			originalValidateS3Existence := validateS3Existence
			originalPushToSQSFn := pushToSQSFn

			if tc.mockValidateFunc != nil {
				validateS3Existence = tc.mockValidateFunc
			}
			if tc.mockPushFunc != nil {
				pushToSQSFn = tc.mockPushFunc
			}
			defer func() {
				validateS3Existence = originalValidateS3Existence
				pushToSQSFn = originalPushToSQSFn
			}()

			err := ValidateS3Path_PushToSQS(&BulkOrderRequest{FilePath: tc.filePath})
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
