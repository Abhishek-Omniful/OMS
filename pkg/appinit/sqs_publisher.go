package appinit

import (
	"github.com/Abhishek-Omniful/OMS/mycontext"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/sqs"
)

var publisher *sqs.Publisher

func PublisherInit(newQueue *sqs.Queue) {
	ctx := mycontext.GetContext()

	logger.Infof(i18n.Translate(ctx, "Initializing SQS Publisher"))
	publisher = sqs.NewPublisher(newQueue)
	logger.Infof(i18n.Translate(ctx, "SQS Publisher successfully created"))
}

func GetPublisher() *sqs.Publisher {
	return publisher
}
