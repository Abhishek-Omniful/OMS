package controllers

import (
	"github.com/Abhishek-Omniful/OMS/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/i18n"
)

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
