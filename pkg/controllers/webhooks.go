package controllers

import (
	"github.com/Abhishek-Omniful/OMS/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"
)

var logger = log.DefaultLogger()

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
