package entities

type WebhookPayload struct {
	AspectType     string         `json:"aspect_type"`
	EventTime      int64          `json:"event_time"`
	ObjectID       int64          `json:"object_id"`
	ObjectType     string         `json:"object_type"`
	OwnerID        int64          `json:"owner_id"`
	SubscriptionID int64          `json:"subscription_id"`
	Updates        map[string]any `json:"updates,omitempty"`
}
