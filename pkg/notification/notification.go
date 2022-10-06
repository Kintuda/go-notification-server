package notification

import "time"

type InitialWebhook struct {
	Payload    map[string]interface{} `json:"payload" binding:"required"`
	NotifierID string                 `json:"notifier_id" binding:"required"`
}

type Webhook struct {
	ID          string                 `json:"id"`
	Payload     map[string]interface{} `json:"payload"`
	NotifierID  string                 `json:"notifier_id"`
	MaxAttempts int                    `json:"max_attempts"`
	CreatedAt   time.Time              `json:"created_at"`
}

type WebhookResponse struct {
	Body   string `json:"body"`
	Status int    `json:"status"`
}
