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
		log.Printf("Falha ao enviar payload do webhook para Kafka: %v", err)
		return err
	}

	log.Printf("Payload do webhook processado: %+v", payload)
	return nil
}
