package http

import (
	"github.com/Kintuda/notification-server/pkg/notification"
	"github.com/gin-gonic/gin"
)

type CreateNotifier struct {
	Endpoint       string `json:"endpoint"`
	EndpointMethod string `json:"endpoint_method"`
	ClientID       string `json:"client_id"`
	ClientSecret   string `json:"client_secret"`
	MaxAttempts    int    `json:"max_attempts"`
}

func (r *Router) RegisterNotifier(c *gin.Context) {
	payload := &notification.CreateNotifier{}

	if err := c.ShouldBindJSON(payload); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	notifier, err := r.NotifierService.CreateNotifier(c.Request.Context(), *payload)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(201, notifier)
}
