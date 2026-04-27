package services

import (
	"log"

	entities "strava-webhook-api/cmd/entities"
	infra "strava-webhook-api/cmd/infra"
)

type WebhookService struct{}

func NewWebhookService() *WebhookService {
	return &WebhookService{}
}

func (ws *WebhookService) ProcessWebhookPayload(payload entities.WebhookPayload) error {
	simplifiedPayload := map[string]any{
		"object_id": payload.ObjectID,
		"owner_id":  payload.OwnerID,
	}

	err := infra.SendMessageJSON(simplifiedPayload)
	if err != nil {
		log.Printf("Failed to send webhook payload to Kafka: %v", err)
		return err
	}

	log.Printf("Processed webhook payload: %+v", payload)
	return nil
}
