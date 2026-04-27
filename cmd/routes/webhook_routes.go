package routes

import (
	webhookcontrollers "strava-webhook-api/cmd/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterWebhookRoutes(r *gin.Engine) {
	controller := webhookcontrollers.NewWebhookController()

	r.GET("/webhook/strava", controller.ControllerWebhookSubscriber)
	r.POST("/webhook/strava", controller.ControllerWebhookReceiver)
}
