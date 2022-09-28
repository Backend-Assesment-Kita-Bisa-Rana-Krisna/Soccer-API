package main

import (
	"os"
	"soccer-api/configuration"
	"soccer-api/route"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("APP_PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"})

	configuration.ConnectDB()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"data": "Hello World",
		})
	})

	api := router.Group("/api")

	route.ApiRoutes(api)

	router.Run(":" + port)

	defer configuration.CloseConnectDB()
}
