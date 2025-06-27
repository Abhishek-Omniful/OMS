package appinit

import (
	"context"
	"strconv"
	"time"

	"github.com/Abhishek-Omniful/OMS/mycontext"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"
	"go.mongodb.org/mongo-driver/bson"
)

type Webhook struct {
	URL      string `json:"url" bson:"url"`
	TenantID int64  `json:"tenant_id" bson:"tenant_id"`
}

var ctx context.Context

func CacheWebhookURL(tenantID int64, url string) {
	log.Infof(i18n.Translate(ctx, "Caching webhook URL for TenantID=%d"), tenantID)
	key := "webhook:" + strconv.FormatInt(tenantID, 10)
	_, err := RedisClient.Set(ctx, key, url, 0)
	if err != nil {
		log.Error(i18n.Translate(ctx, "Failed to cache webhook URL"), err)
	}
	log.Infof(i18n.Translate(ctx, "Cached webhook URL for TenantID=%d"), tenantID)
}

func CheckCache(tenantID int64) string {
	log.Infof(i18n.Translate(ctx, "Checking cache for webhook URL for TenantID=%d"), tenantID)
	key := "webhook:" + strconv.FormatInt(tenantID, 10)
	val, err := RedisClient.Get(ctx, key)
	if err != nil {
		log.Warn(i18n.Translate(ctx, "Failed to get webhook URL from cache"), err)
		return ""
	}
	return val
}

func PostToWebhook(tenantID int64, urlStr string, payload interface{}) {
	request := &http.Request{
		Url: urlStr,
		Body: map[string]interface{}{
			"body": payload,
		},
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Timeout: 5 * time.Second, // Optional request-specific timeout
	}
	_, err := client.Post(request, nil)
	if err != nil {
		logger.Error(i18n.Translate(ctx, "Failed to send POST request to webhook URL"), err)
		return
	}
}

func SendNotification(tenantID int64, payload interface{}) {
	ctx = mycontext.GetContext()
	urlStr := CheckCache(tenantID)
	if urlStr == "" {
		var wh Webhook
		err := WebhookCollection.FindOne(ctx, bson.M{"tenant_id": tenantID}).Decode(&wh)
		log.Println(tenantID)
		if err != nil {
			log.Warn(i18n.Translate(ctx, "No webhook found for tenant"), tenantID)
			return
		}
		urlStr = wh.URL
		CacheWebhookURL(tenantID, urlStr)
	}

	PostToWebhook(tenantID, urlStr, payload)

	log.Infof(i18n.Translate(ctx, "Successfully sent webhook for TenantID=%d to URL=%s"), tenantID, urlStr)
}
