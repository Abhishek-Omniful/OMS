package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/i18n"
)

func ServeHome(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": i18n.Translate(c, "OMS Service"),
	})
}
