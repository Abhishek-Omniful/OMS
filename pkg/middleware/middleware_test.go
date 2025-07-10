package middlewares

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLogRequestMiddleware_UnitStyle(t *testing.T) {
	router := gin.New()
	ctx := context.Background()

	// Attach middleware and register route and handler
	router.Use(LogRequest(ctx))
	router.GET("/test", func(c *gin.Context) {
		time.Sleep(5 * time.Millisecond) // simulate latency
		c.String(http.StatusOK, "ok")
	})

	// Set up request with headers
	req := httptest.NewRequest("GET", "/test?foo=bar", bytes.NewBuffer(nil))
	req.Header.Set("X-Tenant-ID", "tenant-123")
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)
	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ok", w.Body.String())
}
