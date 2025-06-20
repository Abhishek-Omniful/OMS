package appinit

import (
	"log"

	"github.com/Abhishek-Omniful/OMS/mycontext"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/sqs"
)

func queueInit() *sqs.Queue {

	log.Println("Initializing SQS Queue")

	ctx := mycontext.GetContext()

	region := config.GetString(ctx, "aws.region")
	account := config.GetString(ctx, "aws.account")
	endpoint := config.GetString(ctx, "aws.sqsendpoint")

	log.Println("Region:", region, "Account:", account, "Endpoint:", endpoint)

	sqsConfig := sqs.GetSQSConfig(
		ctx,
		true,
		"",
		region,
		account,
		endpoint,
	)

	log.Println("SQS Config:", sqsConfig)

	newQueue, err := sqs.NewStandardQueue(ctx, "sqs-queue", sqsConfig)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Standard SQS Queue Successfully created")
	return newQueue
}

func GetSqs() *sqs.Queue {
	return queueInit()
}
