package router

import (
	"backend-service/internal/adapter/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(orderHandler handler.OrderHandlerInterface) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.POST("/orders", orderHandler.CreateOrder)

	return r
}
