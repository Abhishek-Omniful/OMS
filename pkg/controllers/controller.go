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
// 	err = models.ValidateS3Path_PushToSQS(req)

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

// }

// func CreateWebhook(c *gin.Context) {
// 	var req = &models.Webhook{}
// 	err := c.ShouldBindBodyWithJSON(&req)
// 	if err != nil {
// 		c.JSON(400, gin.H{
// 			"error": "Failed Parse request",
// 		})
// 		return
// 	}
// 	err = models.CreateWebhook(req)
// 	if err != nil {
// 		c.JSON(400, gin.H{
// 			"error": "Failed to create webhook",
// 		})
// 		return
// 	}
// 	c.JSON(200, gin.H{
// 		"message": "Webhook created successfully!",
// 	})
// }

// func ListWebhooks(c *gin.Context) {
// 	webhooks, err := models.ListWebhooks()
// 	if err != nil {
// 		c.JSON(500, gin.H{
// 			"error": "Failed to list webhooks",
// 		})
// 		return
// 	}
// 	c.JSON(200, gin.H{
// 		"webhooks": webhooks,
// 	})
// }

package controllers

import (
	"github.com/Abhishek-Omniful/OMS/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"
)

var logger = log.DefaultLogger()

func ServeHome(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": i18n.Translate(c, "OMS Service"),
	})
}

func StoreInS3(c *gin.Context) {
	var req = &models.StoreCSV{}
	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			i18n.Translate(c, "error"): i18n.Translate(c, "Failed Parse request"),
		})
		return
	}
	err = models.StoreInS3(req)
	if err != nil {
		c.JSON(400, gin.H{
			i18n.Translate(c, "error"): i18n.Translate(c, "Failed to upload to s3"),
		})
		return
	}
	c.JSON(200, gin.H{
		i18n.Translate(c, "message"): i18n.Translate(c, "File uploaded to S3!"),
	})
}

func CreateBulkOrder(c *gin.Context) {
	var req = &models.BulkOrderRequest{}
	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			i18n.Translate(c, "error"): i18n.Translate(c, "Failed Parse request"),
		})
		return
	}
	err = models.ValidateS3Path_PushToSQS(req)

	if err != nil {
		c.JSON(400, gin.H{
			i18n.Translate(c, "error"): i18n.Translate(c, "Invalid  path to s3 or s3 bucket dont exits, first try creatring one and retry"),
		})
		return
	}
	logger.Infof(i18n.Translate(c, "Valid Path to s3 !"))
	logger.Infof(i18n.Translate(c, "Pushing to sqs !"))

	c.JSON(200, gin.H{
		i18n.Translate(c, "message"): i18n.Translate(c, "Valid Path to s3 !"),
	})
}

func CreateWebhook(c *gin.Context) {
	var req = &models.Webhook{}
	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			i18n.Translate(c, "error"): i18n.Translate(c, "Failed Parse request"),
		})
		return
	}
	err = models.CreateWebhook(req)
	if err != nil {
		c.JSON(400, gin.H{
			i18n.Translate(c, "error"): i18n.Translate(c, "Failed to create webhook"),
		})
		return
	}
	c.JSON(200, gin.H{
		i18n.Translate(c, "message"): i18n.Translate(c, "Webhook created successfully!"),
	})
}

func ListWebhooks(c *gin.Context) {
	webhooks, err := models.ListWebhooks()
	if err != nil {
		c.JSON(500, gin.H{
			i18n.Translate(c, "error"): i18n.Translate(c, "Failed to list webhooks"),
		})
		return
	}
	c.JSON(200, gin.H{
		i18n.Translate(c, "webhooks"): webhooks,
	})
}
