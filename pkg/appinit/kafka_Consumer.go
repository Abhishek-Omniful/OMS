package appinit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/kafka"
	"github.com/omniful/go_commons/pubsub"
	"github.com/omniful/go_commons/pubsub/interceptor"
)

var kafkaConsumer *kafka.ConsumerClient

func checkInventory(sku_id, hub_id int64, quantity int) bool {
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
		logger.Errorf("Failed to call IMS validate API: %v", err)
		return false
	}
	logger.Infof("Response from IMS validate API: %v", response)
	return response.IsValid
}

// Implement message handler
type MessageHandler struct{}

func (h *MessageHandler) Handle(ctx context.Context, msg *pubsub.Message) error {

	log.Printf("Handling message from topic: %s, key: %s, value: %s", msg.Topic, msg.Key, msg.Value)

	var order Order
	if err := json.Unmarshal(msg.Value, &order); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return err
	}
	sku_id := order.SKUID
	hub_id := order.HubID
	quantity := order.Quantity

	status := checkInventory(sku_id, hub_id, quantity)
	if !status {
		logger.Println("Not Enough Inventory ! , Keep it on Hold")
		return nil
	}

	logger.Printf("Processing order: %+v", order)

	// changeThe status of order to ""new Order" in mongoDB
	order.Status = "new Order"
	err = SaveOrder(ctx, &order, ordersCollection)
	if err != nil {
		logger.Printf("Failed to save order: %v", err)
		return err
	}
	return nil
}

// Implement the required Process method for IPubSubMessageHandler interface
func (h *MessageHandler) Process(ctx context.Context, msg *pubsub.Message) error {
	return h.Handle(ctx, msg)
}

func InitKafkaConsumer() {
	log.Println("Initializing Kafka consumer...")

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
	// defer func() {
	// 	log.Println("Closing Kafka consumer")
	// 	kafkaConsumer.Close()
	// }()

	log.Println("Attaching NewRelic interceptor to consumer")
	kafkaConsumer.SetInterceptor(interceptor.NewRelicInterceptor())

	handler := &MessageHandler{}
	topic := "my-topic"

	log.Printf("Registering handler for topic: %s", topic)
	kafkaConsumer.RegisterHandler(topic, handler)

	log.Printf("Subscribing to topic: %s", topic)
	ctx := context.Background()
	go kafkaConsumer.Subscribe(ctx) // running as a goroutine

	// BLOCK forever so consumer can keep running
	select {}
}
