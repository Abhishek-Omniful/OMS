package dbService

import (
	"github.com/Abhishek-Omniful/OMS/mycontext"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var logger = log.DefaultLogger()

var mongoClient *mongo.Client
var err error

var (
	OrdersCollection  *mongo.Collection
	WebhookCollection *mongo.Collection
)

func ConnectDB() {
	ctx := mycontext.GetContext()
	logger.Infof(i18n.Translate(ctx, "Connecting to MongoDB"))

	mongoURI := config.GetString(ctx, "mongo.uri")
	logger.Infof(i18n.Translate(ctx, "MongoDB URI: %s"), mongoURI)

	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		logger.Error(i18n.Translate(ctx, "Failed to connect to MongoDB"), err)
		return
	}
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		logger.Error(i18n.Translate(ctx, "Failed to ping MongoDB"), err)
		return
	}
	logger.Infof(i18n.Translate(ctx, "Connected to MongoDB successfully"))

	dbName := config.GetString(ctx, "mongo.dbname")
	webhookCollectionName := config.GetString(ctx, "mongo.webhooksCollection")
	ordersCollectionName := config.GetString(ctx, "mongo.ordersCollection")

	OrdersCollection = mongoClient.Database(dbName).Collection(ordersCollectionName)
	logger.Infof(i18n.Translate(ctx, "Connected to MongoDB collection: %s"), ordersCollectionName)

	WebhookCollection = mongoClient.Database(dbName).Collection(webhookCollectionName)
	logger.Infof(i18n.Translate(ctx, "Connected to MongoDB collection: %s"), webhookCollectionName)
}

func GetDB() *mongo.Client {
	return mongoClient
}

func GetOrdersCollection() *mongo.Collection {
	return OrdersCollection
}

func GetWebhookCollection() *mongo.Collection {
	return WebhookCollection
}
