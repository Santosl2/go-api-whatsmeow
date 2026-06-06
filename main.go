package main

import (
	"fmt"

	instances "github.com/Santosl2/go-api-whatsmeow/pkg/instances/handler"
	"github.com/Santosl2/go-api-whatsmeow/pkg/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mau.fi/whatsmeow"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Error loading .env file")
	}

	server := gin.Default()

	// CORS middleware
	server.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Cache-Control, X-Requested-With, apikey, ApiKey")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	})

	whatsmeowService := services.NewWhatsmeowService(services.SQLStore(), make(map[string]*whatsmeow.Client))

	whatsmeowService.StartAllConnections()

	instancesHandler := instances.NewInstancesHandler(whatsmeowService)

	// Routes
	server.GET("/instances", instancesHandler.GetInstances)
	server.POST("/instances", instancesHandler.CreateInstance)

	server.Run(":3232")
}
