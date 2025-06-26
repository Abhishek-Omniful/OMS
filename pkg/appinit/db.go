// package appinit

// import (
// 	"log"

// 	"github.com/Abhishek-Omniful/OMS/mycontext"
// 	"github.com/omniful/go_commons/config"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// 	"go.mongodb.org/mongo-driver/mongo/readpref"
// )

// var mongoClient *mongo.Client
// var err error
// var collection *mongo.Collection

// func ConnectDB() {
// 	ctx := mycontext.GetContext()
// 	log.Println("Connecting to MongoDB")
// 	mongoURI := config.GetString(ctx, "mongo.uri")
// 	log.Printf("MongoDB URI: %s", mongoURI)
// 	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
// 	if err != nil {
// 		log.Fatal(err)
// 		return
// 	}
// 	err = mongoClient.Ping(ctx, readpref.Primary())
// 	if err != nil {
// 		log.Fatal(err)
// 		return
// 	}
// 	log.Println("Connected to MongoDB successfully")
// }

// func GetDB() *mongo.Client {
// 	return mongoClient
// }

// func GetMongoCollection(dbname string, collectionName string) (*mongo.Collection, error) {
// 	mongoClient := GetDB()
// 	collection = mongoClient.Database(dbname).Collection(collectionName)
// 	log.Printf("Connected to MongoDB collection: %s in database: %s", collectionName, dbname)
// 	return collection, err
// }

package appinit

import (
	"log"

	"github.com/Abhishek-Omniful/OMS/mycontext"
	"github.com/omniful/go_commons/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoClient *mongo.Client
var err error

var (
	ordersCollection  *mongo.Collection
	webhookCollection *mongo.Collection
)

// Connects to MongoDB and initializes both collections
func ConnectDB() {
	ctx := mycontext.GetContext()
	log.Println("Connecting to MongoDB")
	mongoURI := config.GetString(ctx, "mongo.uri")
	log.Printf("MongoDB URI: %s", mongoURI)

	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
		return
	}
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Connected to MongoDB successfully")

	// Initialize both collections
	dbName := config.GetString(ctx, "mongo.dbname")
	webhookCollectionName := config.GetString(ctx, "mongo.webhooksCollection")
	ordersCollectionName := config.GetString(ctx, "mongo.ordersCollection")
	ordersCollection = mongoClient.Database(dbName).Collection(ordersCollectionName)
	logger.Infof("Connected to MongoDB collection: %s", ordersCollectionName)
	webhookCollection = mongoClient.Database(dbName).Collection(webhookCollectionName)
	logger.Infof("Connected to MongoDB collection: %s", webhookCollectionName)
}

// Returns the raw Mongo client
func GetDB() *mongo.Client {
	return mongoClient
}

// Accessor for the orders collection
func GetOrdersCollection() *mongo.Collection {
	return ordersCollection
}

// Accessor for the webhook collection
func GetWebhookCollection() *mongo.Collection {
	return webhookCollection
}
