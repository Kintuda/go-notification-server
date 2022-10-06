package notification

import "time"

type Notifier struct {
	ID             string    `json:"id"`
	Endpoint       string    `json:"endpoint"`
	EndpointMethod string    `json:"endpoint_method"`
	ClientID       *string   `json:"client_id"`
	ClientSecret   *string   `json:"client_secret"`
	MaxAttempts    int       `json:"max_attempts"`
	CreatedAt      time.Time `json:"created_at"`
}
