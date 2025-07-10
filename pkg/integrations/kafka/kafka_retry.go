package kafkaService

import (
	"context"
	"fmt"
	"time"

	"github.com/Abhishek-Omniful/OMS/mycontext"
	"github.com/Abhishek-Omniful/OMS/pkg/helper/common"
	dbService "github.com/Abhishek-Omniful/OMS/pkg/integrations/db"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ordersCollection *mongo.Collection
var logger *log.Logger = log.DefaultLogger()

func OrderRetryWorker() {
	ctx := mycontext.GetContext()
	ordersCollection = dbService.GetOrdersCollection()
	go func() {
		ticker := time.NewTicker(2 * time.Minute) // after every 2 mins
		defer ticker.Stop()

		for range ticker.C {
			logger.Info(i18n.Translate(ctx, "retrying on_hold orders..."))
			ProcessOnHoldOrders()
		}
	}()
}

func ProcessOnHoldOrders() {
	ctx := mycontext.GetContext()

	orders, err := GetOnHoldOrders(ctx)
	if err != nil {
		logger.Errorf(i18n.Translate(ctx, "Failed to fetch on_hold orders: %v"), err)
		return
	}

	for _, order := range orders {
		logger.Infof(i18n.Translate(ctx, "Retrying order: %d"), order.OrderID)
		CheckInventory(order.SKUID, order.HubID, order.Quantity)
	}
}
func GetOnHoldOrders(ctx context.Context) ([]common.Order, error) {
	var orders []common.Order
    
	if ordersCollection == nil {
		return nil, fmt.Errorf("ordersCollection is nil. Ensure MongoDB is properly initialized")
	}
    
	cursor, err := ordersCollection.Find(ctx, bson.M{"status": "on_hold"})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var order = common.Order{}
		if err := cursor.Decode(&order); err != nil {
			continue
		}
		orders = append(orders, order)
	}

	return orders, nil
}
