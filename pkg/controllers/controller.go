package controllers

import (
	"log"
	"net/http"

	"github.com/Abhishek-Omniful/OMS/mycontext"
	"github.com/Abhishek-Omniful/OMS/pkg/appinit"
	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/config"
	"go.mongodb.org/mongo-driver/mongo"
)

// var mongoClinet *mongo.Client
var collection *mongo.Collection
var err error

func init() {
	ctx := mycontext.GetContext()
	dbname := config.GetString(ctx, "mongo.dbname")
	collectionName := config.GetString(ctx, "mongo.collectionName")

	_, err = appinit.GetDB()
	if err != nil {
		log.Fatal(err)
	}
	collection, err = appinit.GetMongoCollection(dbname, collectionName)
	if err != nil {
		log.Fatal(err)
	}
}

func ServeHome(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"messgae": "OMS Service",
	})
}

func CreateBulkOrder(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"messgae": "Bulk order",
	})
}
