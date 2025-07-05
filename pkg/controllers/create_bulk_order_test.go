package controllers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Abhishek-Omniful/OMS/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ---- MOCKING ----
var validateCalled = false

func mockValidateS3Path_PushToSQS(req *models.BulkOrderRequest) error {
	validateCalled = true
	return nil // simulate success
}

// ---- TEST: Success ----
func TestCreateBulkOrder_Success(t *testing.T) {
	// Override the real function with our mock
	originalFunc := models.ValidateS3Path_PushToSQS
	models.ValidateS3Path_PushToSQS = mockValidateS3Path_PushToSQS
	defer func() { models.ValidateS3Path_PushToSQS = originalFunc }() // restore original after test

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/bulk-order", CreateBulkOrder)

	// Valid request body
	jsonPayload := `{"filePath": "s3://my-bucket/my-file.csv"}`
	req := httptest.NewRequest(http.MethodPost, "/bulk-order", bytes.NewBufferString(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
	assert.True(t, validateCalled)
}

// ---- TEST: Bad Request (invalid JSON) ----
func TestCreateBulkOrder_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/bulk-order", CreateBulkOrder)

	// Invalid request body
	jsonPayload := `{"filePath": }`
	req := httptest.NewRequest(http.MethodPost, "/bulk-order", bytes.NewBufferString(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 400, resp.Code)
}
