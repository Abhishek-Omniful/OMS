package helper

import (
	"log"

	"github.com/Abhishek-Omniful/OMS/pkg/appinit"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetMongoCollection(dbname string, collectionName string) (*mongo.Collection, error) {
	mongoClient, err := appinit.GetDB()
	if err != nil {
		log.Fatal(err)
	}
	collection := mongoClient.Database(dbname).Collection(collectionName)
	return collection, err
}
