package appinit

import (
	"log"

	"github.com/omniful/go_commons/sqs"
)

func PublisherInit(newQueue *sqs.Queue) *sqs.Publisher {
	log.Println("Initializing SQS Publisher")
	publisher := sqs.NewPublisher(newQueue)
	log.Println("SQS Publisher successfully created")
	return publisher
}

func GetPublisher(newQueue *sqs.Queue) *sqs.Publisher {
	return PublisherInit(newQueue)
}
