package init

import (
	"log"

	"time"

	"github.com/Abhishek-Omniful/OMS/mycontext"
	"github.com/Abhishek-Omniful/OMS/pkg/services"
)

func init() {
	log.Println("Initializing application...1")
	ctx := mycontext.GetContext()
	services.ConnectDB()
	// ordersCollection := services.GetOrdersCollection()
	// webhookCollection := services.GetWebhookCollection()

	services.ConnectS3()
	// client := services.GetS3Client()

	services.SQSInit()
	newQueue := services.GetSqs()

	services.PublisherInit(newQueue)
	// publisher := services.GetPublisher()

	services.InitConsumer()
	services.StartConsumer(ctx)

	go services.InitKafkaConsumer()
	time.Sleep(3 * time.Second)
	services.InitKafkaProducer()
	services.ConnectRedis()
	services.OrderRetryWorker()
}

func Initialize() {
	log.Println("Application initialized successfully")
}
