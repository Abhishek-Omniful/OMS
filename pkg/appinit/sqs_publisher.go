package appinit

import (
	"log"

	"github.com/omniful/go_commons/sqs"
)

func PublisherInit() *sqs.Publisher {
	log.Println("Initializing SQS Publisher")
	newQueue := GetSqs()
	publisher := sqs.NewPublisher(newQueue)
	log.Println("SQS Publisher successfully created")
	return publisher
}

func GetPublisher() *sqs.Publisher {
	return PublisherInit()
}
