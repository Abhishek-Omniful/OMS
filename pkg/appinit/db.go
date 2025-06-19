package appinit

import (
	"log"

	"github.com/Abhishek-Omniful/OMS/mycontext"
	"github.com/omniful/go_commons/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func connectDB() (*mongo.Client, error) {
	ctx := mycontext.GetContext()
	log.Println("Connecting to MongoDB")
	mongoURI := config.GetString(ctx, "mongo.uri")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	log.Println("Connected to MongoDB successfully")
	return client, nil
}

func GetDB() (*mongo.Client, error) {
	mongoClient, err := connectDB()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return mongoClient, nil
}

func GetMongoCollection(dbname string, collectionName string) (*mongo.Collection, error) {
	mongoClient, err := GetDB()
	if err != nil {
		log.Fatal(err)
	}
	collection := mongoClient.Database(dbname).Collection(collectionName)
	return collection, err
}
