package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Abhishek-Omniful/OMS/pkg/helper/common"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	baseURL        = "http://localhost:8082"
	testOrderID    = int64(1011)
	expectedStatus = "new Order"
)

var logger = log.DefaultLogger()
var mongoClient *mongo.Client
var OrdersCollection *mongo.Collection

func ConnectDB() {
	ctx := context.Background()
	logger.Infof(i18n.Translate(ctx, "Connecting to MongoDB"))

	mongoURI := "mongodb://localhost:27017"
	logger.Infof(i18n.Translate(ctx, "MongoDB URI: %s"), mongoURI)

	var err error
	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		logger.Error(i18n.Translate(ctx, "Failed to connect to MongoDB"), err)
		return
	}
	if err = mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		logger.Error(i18n.Translate(ctx, "Failed to ping MongoDB"), err)
		return
	}
	logger.Infof(i18n.Translate(ctx, "Connected to MongoDB successfully"))

	OrdersCollection = mongoClient.Database("OMS").Collection("orders")
	logger.Infof(i18n.Translate(ctx, "Connected to MongoDB collection: %s"), "orders")
}

func TestEndToEnd_OrderWorkflow(t *testing.T) {
	ConnectDB()
	ctx := context.Background()

	// Call bulkorder API
	t.Run("Submit S3 CSV Path", func(t *testing.T) {
		payload := map[string]string{
			"filePath": "s3://oms-temp-public/orders.csv",
		}
		body, _ := json.Marshal(payload)

		resp, err := http.Post(fmt.Sprintf("%s/api/v1/order/bulkorder", baseURL), "application/json", bytes.NewBuffer(body))
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		//  Call bulkorder API
		t.Run("Submit S3 CSV Path", func(t *testing.T) {
			payload := map[string]string{
				"filePath": "s3://oms-temp-public/orders.csv",
			}
			body, _ := json.Marshal(payload)

			resp, err := http.Post(fmt.Sprintf("%s/api/v1/order/bulkorder", baseURL), "application/json", bytes.NewBuffer(body))
			assert.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode)
		})

	})

	// Wait for status to change to "new Order"
	t.Run("Wait for status change to new Order", func(t *testing.T) {
		assert.Eventually(t, func() bool {
			order, err := fetchOrderByID(ctx, testOrderID)
			if err != nil {
				t.Logf("Error fetching order: %v", err)
				return false
			}
			t.Logf("Current status: %s", order.Status)
			return order.Status == expectedStatus
		}, 10*time.Second, 2*time.Second)
	})
}

func fetchOrderByID(ctx context.Context, orderID int64) (*common.Order, error) {
	fmt.Println("Fetching order by ID:", orderID)
	var order common.Order
	err := OrdersCollection.FindOne(ctx, bson.M{
		"tenant_id": 1,
		"order_id":  orderID,
	}).Decode(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}
