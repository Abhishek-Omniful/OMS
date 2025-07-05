package controllers

import (
	"github.com/Abhishek-Omniful/OMS/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/i18n"
)

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
