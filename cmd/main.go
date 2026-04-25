package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type WebhookPayload struct {
	AspectType     string         `json:"aspect_type"`
	EventTime      int64          `json:"event_time"`
	ObjectID       int64          `json:"object_id"`
	ObjectType     string         `json:"object_type"`
	OwnerID        int64          `json:"owner_id"`
	SubscriptionID int64          `json:"subscription_id"`
	Updates        map[string]any `json:"updates"`
}

type KafkaPayload struct {
	ObjectID int64 `json:"object_id"`
	OwnerID  int64 `json:"owner_id"`
}

func main() {
	r := gin.Default()

	r.GET("/webhook/strava", handleWebhookGet)
	r.POST("/webhook/strava", handleWebhookPost)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on :%s", port)
	log.Fatal(r.Run(":" + port))
}

func handleWebhookGet(c *gin.Context) {
	challenge := c.Query("hub.challenge")

	if challenge == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing challenge",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"hub.challenge": challenge,
	})
}

func handleWebhookPost(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}

	log.Printf("Webhook received: %s", string(body))

	var payload WebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Printf("Invalid JSON: %v", err)
	} else {
		log.Printf("Parsed payload: %+v", payload)
	}

	if payload.ObjectType == "activity" && payload.AspectType == "create" {
		kafkaPayload := KafkaPayload{
			ObjectID: payload.ObjectID,
			OwnerID:  payload.OwnerID,
		}

		log.Printf("Would send to Kafka: %+v", kafkaPayload)
	}

	c.Status(http.StatusOK)
}
