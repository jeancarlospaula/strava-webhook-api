package main

import (
	"log"
	"net/http"
	"os"

	infra "strava-webhook-api/cmd/infra"
	routes "strava-webhook-api/cmd/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Warning: .env not loaded")
		}
	}

	infra.ConnectKafka()

	r := gin.Default()
	routes.RegisterWebhookRoutes(r)

	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on :%s", port)
	log.Fatal(r.Run(":" + port))
}
