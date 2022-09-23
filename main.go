package main

import (
	"soccer-api/configuration"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"})

	configuration.ConnectDB()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"data": "Hello World",
		})
	})

	router.Run("localhost:3000")

	defer configuration.CloseConnectDB()
}
