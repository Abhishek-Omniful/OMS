package appinit

import (
	"github.com/Abhishek-Omniful/OMS/mycontext"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/sqs"
)

var newQueue *sqs.Queue

func SQSInit() {
	ctx := mycontext.GetContext()

	logger.Infof(i18n.Translate(ctx, "Initializing SQS Queue"))

	region := config.GetString(ctx, "aws.region")
	account := config.GetString(ctx, "aws.account")
	endpoint := config.GetString(ctx, "aws.sqsendpoint")
	queueName := config.GetString(ctx, "aws.sqsname")

	sqsCfg := sqs.GetSQSConfig(
		ctx,
		false,
		"",
		region,
		account,
		endpoint,
	)

	logger.Infof(i18n.Translate(ctx, "SQS Config initialized: %+v"), sqsCfg)

	// Ensure queue exists (create if missing)
	err := sqs.CreateQueue(ctx, sqsCfg, queueName, "standard")
	if err != nil {
		logger.Panicf(i18n.Translate(ctx, "Error creating SQS queue: %v"), err)
		return
	}

	newQueue, err = sqs.NewStandardQueue(ctx, queueName, sqsCfg)
	if err != nil {
		logger.Panicf(i18n.Translate(ctx, "Failed to initialize SQS queue: %v"), err)
		return
	}

	logger.Infof(i18n.Translate(ctx, "Standard SQS Queue successfully created"))
}

func GetSqs() *sqs.Queue {
	return newQueue
}
