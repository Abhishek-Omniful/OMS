package init

import (
	"log"

	"time"

	"github.com/Abhishek-Omniful/OMS/mycontext"
	dbService "github.com/Abhishek-Omniful/OMS/pkg/integrations/db"
	httpclient "github.com/Abhishek-Omniful/OMS/pkg/integrations/httpClient"
	kafkaService "github.com/Abhishek-Omniful/OMS/pkg/integrations/kafka"
	redisService "github.com/Abhishek-Omniful/OMS/pkg/integrations/redis"
	s3Service "github.com/Abhishek-Omniful/OMS/pkg/integrations/s3"
	sqsService "github.com/Abhishek-Omniful/OMS/pkg/integrations/sqs"
)

func init() {
	ctx := mycontext.GetContext()
	dbService.ConnectDB()
	httpclient.InitHttpClient()
	s3Service.ConnectS3()

	go kafkaService.InitKafkaConsumer()
	time.Sleep(3 * time.Second)
	kafkaService.InitKafkaProducer()
	kafkaService.OrderRetryWorker()
	
	redisService.ConnectRedis()
	sqsService.SQSInit()
	newQueue := sqsService.GetSqs()
	sqsService.InitSqsPublisher(newQueue)
	sqsService.InitSqsConsumer()
	sqsService.StartConsumer(ctx)

}

func Initialize() {
	log.Println("Application initialized successfully")
}
