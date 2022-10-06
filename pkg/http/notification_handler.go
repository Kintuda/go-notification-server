package http

import (
	"github.com/Kintuda/notification-server/pkg/notification"
	"github.com/gin-gonic/gin"
)

func (r *Router) SendNotification(c *gin.Context) {
	payload := &notification.InitialWebhook{}

	if err := c.ShouldBindJSON(payload); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	webhook, err := r.NotificationService.SendWebhook(c.Request.Context(), *payload)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(201, webhook)
}
