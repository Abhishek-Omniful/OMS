package appinit

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/omniful/go_commons/kafka"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/pubsub"
)

var kafkaProducer *kafka.ProducerClient

func InitKafkaProducer() {
	log.Println("Initializing Kafka producer")

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
	
	// Marshal order into JSON
	jsonBytes, err := json.Marshal(order)
	if err != nil {
		log.Errorf("Failed to marshal order: %v", err)
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

	ctx := context.WithValue(context.Background(), "request_id", fmt.Sprintf("req-%d", order.OrderID))

	log.Infof("Publishing order to topic: %s", msg.Topic)
	err = kafkaProducer.Publish(ctx, msg)
	if err != nil {
		log.Errorf("Failed to publish order: %v", err)
		panic(err)
	}
	log.Infof("Order published to Kafka successfully: OrderID=%d", order.OrderID)
}
