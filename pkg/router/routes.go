package router

import (
	"github.com/Abhishek-Omniful/OMS/pkg/controllers"
	"github.com/omniful/go_commons/http"
)

func Initialize(s *http.Server) {

	s.GET("/", controllers.ServeHome)

	v1 := s.Engine.Group("/api/v1")
	{
		orders := v1.Group("/order") //containing csv file path
		{
			orders.POST("/bulkorder", controllers.CreateBulkOrder)
		}

		csv := v1.Group("/csv")
		{
			csv.POST("/filepath", controllers.StoreInS3)
		}

		webhooks := v1.Group("/webhook")
		{
			webhooks.POST("/create", controllers.CreateWebhook)
			webhooks.GET("/list", controllers.ListWebhooks)
		}
	}
}


