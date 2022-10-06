package notification

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CreateNotifier struct {
	Endpoint       string  `json:"endpoint" binding:"required"`
	EndpointMethod string  `json:"endpoint_method" binding:"required"`
	ClientID       *string `json:"client_id"`
	ClientSecret   *string `json:"client_secret"`
	MaxAttempts    int     `json:"max_attempts" binding:"required"`
}

type NotifierService struct {
	NotificationRepository NotificationRepository
}

func NewNotifierService(repo NotificationRepository) *NotifierService {
	return &NotifierService{NotificationRepository: repo}
}

func (n *NotifierService) CreateNotifier(ctx context.Context, payload CreateNotifier) (*Notifier, error) {
	notifier := Notifier{
		ID:             uuid.NewString(),
		Endpoint:       payload.Endpoint,
		EndpointMethod: strings.ToUpper(payload.EndpointMethod),
		ClientID:       payload.ClientID,
		MaxAttempts:    payload.MaxAttempts,
		ClientSecret:   payload.ClientSecret,
		CreatedAt:      time.Now(),
	}

	if err := n.NotificationRepository.RegisterNotifier(ctx, notifier); err != nil {
		return nil, err
	}

	return &notifier, nil
}
