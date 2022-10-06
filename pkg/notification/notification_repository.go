package notification

import "context"

type Event struct {
	ID       string                 `json:"id"`
	Url      string                 `json:"url"`
	Payload  map[string]interface{} `json:"payload"`
	Attempts int                    `json:"max_attempts"`
}

type NotificationRepository interface {
	FindAttempts(ctx context.Context, id string) (int, error)
	CreateWebhookTransaction(ctx context.Context, w Webhook) error
	RegisterNotifier(ctx context.Context, notifier Notifier) error
	FetchNotifier(ctx context.Context, id string) (*Notifier, error)
}
