package models

import (
	"log"

	"github.com/Abhishek-Omniful/OMS/pkg/appinit"
	"go.mongodb.org/mongo-driver/mongo"
)

type BulkOrderRequest struct {
	SellerID string `json:"sellerID"`
	FilePath string `json:"filePath"`
}

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
