package init

import (
	"log"

	"time"

	"github.com/Abhishek-Omniful/OMS/mycontext"
	"github.com/Abhishek-Omniful/OMS/pkg/services"
)

func init() {
	ctx := mycontext.GetContext()
	services.ConnectDB()
	services.ConnectS3()
	services.SQSInit()
	newQueue := services.GetSqs()
	services.PublisherInit(newQueue)
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
