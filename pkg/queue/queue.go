package queue

import "context"

type QueueProvider interface {
	Publish(ctx context.Context, contentType string, payload []byte) error
}
