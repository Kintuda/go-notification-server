package postgres

import (
	"context"

	"github.com/Kintuda/notification-server/pkg/notification"
	"github.com/georgysavva/scany/pgxscan"
)

type NotificationRepositoryPg struct {
	Executor QueryExecutor
}

func (n *NotificationRepositoryPg) FindAttempts(ctx context.Context, id string) (int, error) {
	return 0, nil
}

func (n *NotificationRepositoryPg) CreateWebhookTransaction(ctx context.Context, w notification.Webhook) error {
	return nil
}

func (n *NotificationRepositoryPg) RegisterNotifier(ctx context.Context, notifier notification.Notifier) error {
	sql := `
	INSERT INTO notifiers (
		id,
		endpoint,
		endpoint_method,
		client_id,
		client_secret,
		max_attempts,
		created_at
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7
	);
`

	_, err := n.Executor.Exec(
		ctx,
		sql,
		notifier.ID,
		notifier.Endpoint,
		notifier.EndpointMethod,
		notifier.ClientID,
		notifier.ClientSecret,
		notifier.MaxAttempts,
		notifier.CreatedAt,
	)

	return err
}

func (n *NotificationRepositoryPg) FetchNotifier(ctx context.Context, id string) (*notification.Notifier, error) {
	var ns []*notification.Notifier
	query := `SELECT * FROM notifiers WHERE id = $1 LIMIT 1`

	if err := pgxscan.Select(ctx, n.Executor, &ns, query, id); err != nil {
		return nil, err
	}

	if len(ns) == 0 {
		return nil, nil
	}

	return ns[0], nil
}
