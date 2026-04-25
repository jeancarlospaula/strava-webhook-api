package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"

	"github.com/joho/godotenv"
)

var producer *kafka.Writer

type WebhookPayload struct {
	AspectType     string         `json:"aspect_type"`
	EventTime      int64          `json:"event_time"`
	ObjectID       int64          `json:"object_id"`
	ObjectType     string         `json:"object_type"`
	OwnerID        int64          `json:"owner_id"`
	SubscriptionID int64          `json:"subscription_id"`
	Updates        map[string]any `json:"updates,omitempty"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env not loaded")
	}

	connectKafka()

	r := gin.Default()

	r.GET("/webhook/strava", ControllerWebhookSubscriber)
	r.POST("/webhook/strava", ControllerWebhookReceiver)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on :%s", port)
	log.Fatal(r.Run(":" + port))
}

func ControllerWebhookSubscriber(c *gin.Context) {
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

func ControllerWebhookReceiver(c *gin.Context) {
	var payload WebhookPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Printf("Invalid JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	log.Printf("Parsed payload: %+v", payload)

	if payload.ObjectType == "activity" && payload.AspectType == "create" {
		sendMessage(payload)
	}

	c.Status(http.StatusOK)
}

func connectKafka() *kafka.Writer {
	CA_CERT := os.Getenv("CA_CERT")
	KAFKA_TOPIC := os.Getenv("KAFKA_TOPIC")
	KAFKA_USERNAME := os.Getenv("KAFKA_USERNAME")
	KAFKA_PASSWORD := os.Getenv("KAFKA_PASSWORD")
	KAFKA_BROKER := os.Getenv("KAFKA_BROKER")

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM([]byte(CA_CERT))
	if !ok {
		log.Fatalf("Failed to parse CA certificate")
	}

	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}

	scram, err := scram.Mechanism(scram.SHA512, KAFKA_USERNAME, KAFKA_PASSWORD)
	if err != nil {
		log.Fatalf("Failed to create scram mechanism: %s", err)
	}

	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		TLS:           tlsConfig,
		SASLMechanism: scram,
	}

	producer = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{KAFKA_BROKER},
		Topic:    KAFKA_TOPIC,
		Balancer: &kafka.Hash{},
		Dialer:   dialer,
	})

	return producer
}

func sendMessage(payload WebhookPayload) {
	message := kafka.Message{Value: fmt.Appendf(nil, `{"object_id": %d, "owner_id": %d}`, payload.ObjectID, payload.OwnerID)}

	err := producer.WriteMessages(context.Background(), message)

	if err != nil {
		log.Printf("failed to write message: %s", err)
	} else {
		log.Printf("message sent to Kafka topic: %s", string(message.Value))
	}
}
