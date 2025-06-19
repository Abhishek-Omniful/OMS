package router

import (
	"github.com/omniful/go_commons/http"
	"github.com/Abhishek-Omniful/OMS/pkg/controllers"
)

func Initialize(s *http.Server) {

	s.GET("/", controllers.ServeHome)

	v1 := s.Engine.Group("/api/v1")
	{
		orders := v1.Group("/order")
		{
			orders.POST("/bulkorder", controllers.CreateBulkOrder)
		}
	}
}
