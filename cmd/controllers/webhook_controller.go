package webhookcontrollers

import (
	"net/http"

	entities "strava-webhook-api/cmd/entities"
	services "strava-webhook-api/cmd/services"

	"github.com/gin-gonic/gin"
)

type webhookController struct {
	webhookService *services.WebhookService
}

func NewWebhookController() *webhookController {
	return &webhookController{
		webhookService: services.NewWebhookService(),
	}
}

func (wc *webhookController) ControllerWebhookSubscriber(c *gin.Context) {
	challenge := c.Query("hub.challenge")

	if challenge == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing challenge",
		})
		return
	}

	c.String(http.StatusOK, challenge)
}

func (wc *webhookController) ControllerWebhookReceiver(c *gin.Context) {
	var payload entities.WebhookPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid payload",
		})
		return
	}

	err := wc.webhookService.ProcessWebhookPayload(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to process webhook",
		})
		return
	}

	c.Status(http.StatusOK)
}
