package appinit

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/kafka"
	"github.com/omniful/go_commons/pubsub"
)

var kafkaProducer *kafka.ProducerClient

func InitKafkaProducer() {
	ctx := context.Background()
	logger.Infof(i18n.Translate(ctx, "Initializing Kafka producer"))

	kafkaProducer = kafka.NewProducer(
		kafka.WithBrokers([]string{"localhost:9092"}),
		kafka.WithClientID("my-producer"),
		kafka.WithKafkaVersion("3.4.0"),
	)
}

func GetKafkaProducer() *kafka.ProducerClient {
	return kafkaProducer
}

func PublishOrder(order *Order) {
	ctx := context.WithValue(context.Background(), "request_id", fmt.Sprintf("req-%d", order.OrderID))

	jsonBytes, err := json.Marshal(order)
	if err != nil {
		logger.Errorf(i18n.Translate(ctx, "Failed to marshal order: %v"), err)
		return
	}

	msg := &pubsub.Message{
		Topic: "my-topic",
		Key:   fmt.Sprintf("order-%d", order.OrderID),
		Value: jsonBytes,
		Headers: map[string]string{
			"source": "order-service",
		},
	}

	logger.Infof(i18n.Translate(ctx, "Publishing order to topic: %s"), msg.Topic)
	err = kafkaProducer.Publish(ctx, msg)
	if err != nil {
		logger.Errorf(i18n.Translate(ctx, "Failed to publish order: %v"), err)
		panic(err)
	}
	logger.Infof(i18n.Translate(ctx, "Order published to Kafka successfully: OrderID=%d"), order.OrderID)
}
