package models

import (
	"log"

	"github.com/Abhishek-Omniful/OMS/pkg/appinit"
	"go.mongodb.org/mongo-driver/mongo"
)

var mongoClinet *mongo.Client
var err error

func init() {
	mongoClinet, err = appinit.GetDB()
	if err != nil {
		log.Fatal(err)
	}
}

func ServeHome() {

}
