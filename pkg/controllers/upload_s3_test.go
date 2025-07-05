package controllers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Abhishek-Omniful/OMS/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ---- MOCK: Replace actual models.StoreInS3 function ----
var storeCalled = false

func mockStoreInS3_Success(s *models.StoreCSV) error {
	storeCalled = true
	return nil
}

func mockStoreInS3_Fail(s *models.StoreCSV) error {
	storeCalled = true
	return errors.New("upload failed")
}

// ---- TEST: Success ----
func TestStoreInS3_Success(t *testing.T) {
	// Backup and mock
	original := models.StoreInS3
	models.StoreInS3 = mockStoreInS3_Success
	defer func() { models.StoreInS3 = original }()

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/upload-to-s3", StoreInS3)

	payload := `{"filePath": "local/path/to/file.csv"}`
	req := httptest.NewRequest(http.MethodPost, "/upload-to-s3", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
	assert.True(t, storeCalled)
}

// ---- TEST: Invalid JSON ----
func TestStoreInS3_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/upload-to-s3", StoreInS3)

	payload := `{"filePath": }` // malformed JSON
	req := httptest.NewRequest(http.MethodPost, "/upload-to-s3", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 400, resp.Code)
}

// ---- TEST: Upload Fail ----
func TestStoreInS3_UploadFails(t *testing.T) {
	// Backup and mock
	original := models.StoreInS3
	models.StoreInS3 = mockStoreInS3_Fail
	defer func() { models.StoreInS3 = original }()

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/upload-to-s3", StoreInS3)

	payload := `{"filePath": "local/path/to/file.csv"}`
	req := httptest.NewRequest(http.MethodPost, "/upload-to-s3", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 400, resp.Code)
	assert.True(t, storeCalled)
}
