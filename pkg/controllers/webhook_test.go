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

// ---------- Mocks ----------
var webhookCreatedCalled = false
var webhookListedCalled = false

func mockCreateWebhookSuccess(req *models.Webhook) error {
	webhookCreatedCalled = true
	return nil
}

func mockCreateWebhookFail(req *models.Webhook) error {
	webhookCreatedCalled = true
	return errors.New("mock insert error")
}

func mockListWebhooksSuccess() ([]models.Webhook, error) {
	webhookListedCalled = true
	return []models.Webhook{
		{URL: "https://example.com", TenantID: 1},
	}, nil
}

func mockListWebhooksFail() ([]models.Webhook, error) {
	webhookListedCalled = true
	return nil, errors.New("mock list error")
}

// ---------- Test: CreateWebhook ----------
func TestCreateWebhook_Success(t *testing.T) {
	original := models.CreateWebhook
	models.CreateWebhook = mockCreateWebhookSuccess
	defer func() { models.CreateWebhook = original }()

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/webhook", CreateWebhook)

	payload := `{"url": "https://example.com", "tenantID": 1}`
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
	assert.True(t, webhookCreatedCalled)
}

func TestCreateWebhook_BadJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/webhook", CreateWebhook)

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBufferString(`{"url": }`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, 400, resp.Code)
}

func TestCreateWebhook_InsertFails(t *testing.T) {
	original := models.CreateWebhook
	models.CreateWebhook = mockCreateWebhookFail
	defer func() { models.CreateWebhook = original }()

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/webhook", CreateWebhook)

	payload := `{"url": "https://fail.com", "tenantID": 1}`
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, 400, resp.Code)
}

// ---------- Test: ListWebhooks ----------
func TestListWebhooks_Success(t *testing.T) {
	original := models.ListWebhooks
	models.ListWebhooks = mockListWebhooksSuccess
	defer func() { models.ListWebhooks = original }()

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/webhooks", ListWebhooks)

	req := httptest.NewRequest(http.MethodGet, "/webhooks", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
	assert.True(t, webhookListedCalled)
	assert.Contains(t, resp.Body.String(), "example.com")
}

func TestListWebhooks_Fails(t *testing.T) {
	original := models.ListWebhooks
	models.ListWebhooks = mockListWebhooksFail
	defer func() { models.ListWebhooks = original }()

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/webhooks", ListWebhooks)

	req := httptest.NewRequest(http.MethodGet, "/webhooks", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, 500, resp.Code)
}
