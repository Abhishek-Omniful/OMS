package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/kafka"
	"github.com/omniful/go_commons/pubsub"
	"github.com/omniful/go_commons/pubsub/interceptor"
)

var kafkaConsumer *kafka.ConsumerClient

func CheckInventory(sku_id, hub_id int64, quantity int) bool {
	ctx := context.Background()
	req := &http.Request{
		Url: "api/v1/inventory/check",
		QueryParams: url.Values{
			"sku_id":   []string{fmt.Sprintf("%d", sku_id)},
			"hub_id":   []string{fmt.Sprintf("%d", hub_id)},
			"quantity": []string{fmt.Sprintf("%d", quantity)},
		},
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Timeout: 5 * time.Second,
	}
	var response ValidationResponse
	_, err := client.Get(req, &response)
	if err != nil {
		logger.Errorf(i18n.Translate(ctx, "Failed to call IMS validate API: %v"), err)
		return false
	}
	logger.Infof(i18n.Translate(ctx, "Response from IMS validate API: %v"), response)
	return response.IsValid
}

// Implement message handler
type MessageHandler struct{}

func (h *MessageHandler) Handle(ctx context.Context, msg *pubsub.Message) error {
	logger.Infof(i18n.Translate(ctx, "Handling message from topic: %s, key: %s, value: %s"), msg.Topic, msg.Key, msg.Value)

	var order Order
	if err := json.Unmarshal(msg.Value, &order); err != nil {
		logger.Errorf(i18n.Translate(ctx, "Failed to unmarshal message: %v"), err)
		return err
	}

	sku_id := order.SKUID
	hub_id := order.HubID
	quantity := order.Quantity

	status := CheckInventory(sku_id, hub_id, quantity)
	if !status {
		logger.Warnf(i18n.Translate(ctx, "Not Enough Inventory! Keeping order on hold"))
		return nil
	}

	logger.Infof(i18n.Translate(ctx, "Processing order: %+v"), order)

	order.Status = "new Order"
	err := SaveOrder(ctx, &order, OrdersCollection)

	logger.Infof(i18n.Translate(ctx, "Notifying the tenant about order creation for TenantID=%d"), order.TenantID)
	SendNotification(order.TenantID, order)

	if err != nil {
		logger.Errorf(i18n.Translate(ctx, "Failed to save order: %v"), err)
		return err
	}
	return nil
}

// IPubSubMessageHandler method
func (h *MessageHandler) Process(ctx context.Context, msg *pubsub.Message) error {
	return h.Handle(ctx, msg)
}

func InitKafkaConsumer() {
	ctx := context.Background()
	logger.Infof(i18n.Translate(ctx, "Initializing Kafka consumer..."))

	kafkaConsumer = kafka.NewConsumer(
		kafka.WithBrokers([]string{"localhost:9092"}),
		kafka.WithConsumerGroup("my-consumer-group"),
		kafka.WithClientID("my-consumer"),
		kafka.WithKafkaVersion("3.4.0"),
	)
	ReceiveOrder()
}

func GetKafkaConsumer() *kafka.ConsumerClient {
	return kafkaConsumer
}

func ReceiveOrder() {
	ctx := context.Background()
	logger.Infof(i18n.Translate(ctx, "Attaching NewRelic interceptor to consumer"))
	kafkaConsumer.SetInterceptor(interceptor.NewRelicInterceptor())

	handler := &MessageHandler{}
	topic := "my-topic"

	logger.Infof(i18n.Translate(ctx, "Registering handler for topic: %s"), topic)
	kafkaConsumer.RegisterHandler(topic, handler)

	logger.Infof(i18n.Translate(ctx, "Subscribing to topic: %s"), topic)
	go kafkaConsumer.Subscribe(ctx)

	select {}
}
