// package controllers

// import (
// 	"log"

// 	"github.com/Abhishek-Omniful/OMS/pkg/models"
// 	"github.com/gin-gonic/gin"
// )

// // var mongoClinet *mongo.Client

// func ServeHome(c *gin.Context) {
// 	c.JSON(200, gin.H{
// 		"message": "OMS Service",
// 	})
// }

// func StoreInS3(c *gin.Context) {
// 	var req = &models.StoreCSV{}
// 	err := c.ShouldBindBodyWithJSON(&req)
// 	if err != nil {
// 		c.JSON(400, gin.H{
// 			"error": "Failed Parse request",
// 		})
// 		return
// 	}
// 	err = models.StoreInS3(req)
// 	if err != nil {
// 		c.JSON(400, gin.H{
// 			"error": "Failed to upload to s3",
// 		})
// 		return
// 	}
// 	c.JSON(200, gin.H{
// 		"message": "File uploaded to S3!",
// 	})
// }

// func CreateBulkOrder(c *gin.Context) {
// 	var req = &models.BulkOrderRequest{}
// 	err := c.ShouldBindBodyWithJSON(&req)
// 	if err != nil {
// 		c.JSON(400, gin.H{
// 			"error": "Failed Parse request",
// 		})
// 		return
// 	}
// 	err = models.ValidateS3Path(req)
// 	log.Println(err)
// 	if err != nil {
// 		c.JSON(400, gin.H{
// 			"error": "Invalid  path to s3 or s3 bucket dont exits, first try creatring one and retry",
// 		})
// 		return
// 	}

// 	log.Println("Valid Path to s3 !")
// 	log.Println("Pushing to sqs !")
// 	c.JSON(200, gin.H{
// 		"message": "Valid Path to s3 !",
// 	})

// 	//push message to sqs

// }

// package controllers

// import (
// 	"log"

// 	// "github.com/Abhishek-Omniful/OMS/mycontext"
// 	// "github.com/Abhishek-Omniful/OMS/pkg/appinit"
// 	"github.com/Abhishek-Omniful/OMS/pkg/models"
// 	"github.com/gin-gonic/gin"
// 	// "github.com/omniful/go_commons/config"
// 	// "go.mongodb.org/mongo-driver/mongo"
// )

// // var mongoClinet *mongo.Client

// var err error

// func init() {

// 	// Initialize S3 client

// }

// func ServeHome(c *gin.Context) {
// 	c.JSON(200, gin.H{
// 		"messgae": "OMS Service",
// 	})
// }

// func StoreInS3(c *gin.Context) {
// 	var req = &models.StoreCSV{}
// 	err := c.ShouldBindBodyWithJSON(&req)
// 	if err != nil {
// 		c.JSON(400, gin.H{
// 			"error": "Failed Parse request",
// 		})
// 		return
// 	}
// 	err = models.StoreInS3(req)
// 	if err != nil {
// 		c.JSON(400, gin.H{
// 			"error": "Failed to upload to s3",
// 		})
// 		return
// 	}
// 	c.JSON(200, gin.H{
// 		"message": "File uploaded to S3!",
// 	})
// }

// func CreateBulkOrder(c *gin.Context) {
// 	var req = &models.BulkOrderRequest{}
// 	err := c.ShouldBindBodyWithJSON(&req)
// 	if err != nil {
// 		c.JSON(400, gin.H{
// 			"error": "Failed Parse request",
// 		})
// 		return
// 	}
// 	err = models.ValidateS3Path(req)

// 	if err != nil {
// 		c.JSON(400, gin.H{
// 			"error": "Invalid  path to s3 or s3 bucket dont exits, first try creatring one and retry",
// 		})
// 		return
// 	}

// 	log.Println("Valid Path to s3 !")
// 	log.Println("Pushing to sqs !")
// 	c.JSON(200, gin.H{
// 		"message": "Valid Path to s3 !",
// 	})

// 	//push message to sqs

// }

package controllers

import (
	"log"

	"github.com/Abhishek-Omniful/OMS/pkg/models"
	"github.com/gin-gonic/gin"
)

// var mongoClinet *mongo.Client

func ServeHome(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "OMS Service",
	})
}

func StoreInS3(c *gin.Context) {
	var req = &models.StoreCSV{}
	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed Parse request",
		})
		return
	}
	err = models.StoreInS3(req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to upload to s3",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "File uploaded to S3!",
	})
}

func CreateBulkOrder(c *gin.Context) {
	var req = &models.BulkOrderRequest{}
	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed Parse request",
		})
		return
	}
	err = models.ValidateS3Path_PushToSQS(req)

	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid  path to s3 or s3 bucket dont exits, first try creatring one and retry",
		})
		return
	}

	log.Println("Valid Path to s3 !")
	log.Println("Pushing to sqs !")
	c.JSON(200, gin.H{
		"message": "Valid Path to s3 !",
	})


}
