package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Kintuda/notification-server/pkg/exception"
	"github.com/Kintuda/notification-server/pkg/queue"
	"github.com/google/uuid"
)

type NotificationService struct {
	NotificationRepository NotificationRepository
	QueueProvider          queue.QueueProvider
}

func NewNotificationService(repo NotificationRepository, provider queue.QueueProvider) *NotificationService {
	return &NotificationService{NotificationRepository: repo, QueueProvider: provider}
}

func (n *NotificationService) SendWebhook(ctx context.Context, initial InitialWebhook) (*Webhook, error) {
	notifier, err := n.NotificationRepository.FetchNotifier(ctx, initial.NotifierID)

	if err != nil {
		return nil, err
	}

	if notifier == nil {
		return nil, &exception.ResourceNotFound{Identifier: initial.NotifierID, Resource: "notifier"}
	}

	w := Webhook{
		ID:          uuid.NewString(),
		NotifierID:  notifier.ID,
		Payload:     initial.Payload,
		MaxAttempts: notifier.MaxAttempts,
		CreatedAt:   time.Now(),
	}

	json, err := json.Marshal(w.Payload)

	if err != nil {
		return nil, err
	}

	if _, err := n.DispatchWebhook(ctx, notifier.EndpointMethod, notifier.Endpoint, json); err != nil {
		return nil, err
	}

	if err := n.NotificationRepository.CreateWebhookTransaction(ctx, w); err != nil {
		return nil, err
	}

	return &w, nil
}

func (n *NotificationService) SendWebhookAsynchronous(ctx context.Context, initial InitialWebhook) (*Webhook, error) {
	notifier, err := n.NotificationRepository.FetchNotifier(ctx, initial.NotifierID)

	if err != nil {
		return nil, err
	}

	if notifier == nil {
		return nil, &exception.ResourceNotFound{Identifier: initial.NotifierID, Resource: "notifier"}
	}

	w := Webhook{
		ID:          uuid.NewString(),
		NotifierID:  notifier.ID,
		Payload:     initial.Payload,
		MaxAttempts: notifier.MaxAttempts,
		CreatedAt:   time.Now(),
	}

	if err != nil {
		return nil, err
	}

	if err := n.NotificationRepository.CreateWebhookTransaction(ctx, w); err != nil {
		return nil, err
	}

	u, err := json.Marshal(w)

	if err != nil {
		return nil, err
	}

	fmt.Println(u)

	if err := n.QueueProvider.Publish(ctx, "text/plain", u); err != nil {
		return nil, err
	}

	return &w, nil
}

func (n *NotificationService) DispatchWebhook(ctx context.Context, method string, URL string, payload []byte) (*WebhookResponse, error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, URL, bytes.NewBuffer(payload))

	if err != nil {
		return nil, err
	}

	response, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	bodyString := string(body)

	if err != nil {
		return nil, err
	}

	return &WebhookResponse{Body: bodyString, Status: response.StatusCode}, nil
}
