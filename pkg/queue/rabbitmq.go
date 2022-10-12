package queue

import (
	"context"
	"fmt"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type RabbitMQProvider struct {
	Conn         *amqp.Connection
	Channel      *amqp.Channel
	Name         string
	QueueName    string
	Done         chan error
	ExchangeName string
	RoutingKey   string
}

type RabbitMQParams struct {
	AmqpDns        string
	Name           string
	QueueName      string
	ExchangeName   string
	RoutingKey     string
	CreateExchange bool
	CreateQueue    bool
}

func NewRabbitMQ(params RabbitMQParams) (*RabbitMQProvider, error) {
	conn, err := amqp.Dial(params.AmqpDns)

	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()

	if err != nil {
		return nil, err
	}

	if params.CreateExchange {
		logrus.Infof("creating exchange: %s", params.ExchangeName)

		if err := channel.ExchangeDeclare(params.ExchangeName, "direct", true, false, false, false, nil); err != nil {
			return nil, err
		}
	}

	if params.CreateQueue {
		logrus.Infof("creating queue: %s", params.QueueName)

		if _, err := channel.QueueDeclare(params.QueueName, true, false, false, false, nil); err != nil {
			return nil, err
		}
	}

	routingKeys := strings.Split(params.RoutingKey, ",")

	for _, routingKey := range routingKeys {
		logrus.Infof("binding queue to consumer routing key, %s", routingKey)

		if err := channel.QueueBind(params.QueueName, routingKey, params.ExchangeName, false, nil); err != nil {
			return nil, err
		}
	}

	return &RabbitMQProvider{
		Conn:         conn,
		Channel:      channel,
		Name:         params.Name,
		QueueName:    params.QueueName,
		ExchangeName: params.ExchangeName,
		// RoutingKey: cfg.RoutingKeys,
	}, nil
}

func (r *RabbitMQProvider) Publish(ctx context.Context, contentType string, payload []byte) error {
	err := r.Channel.PublishWithContext(ctx, r.ExchangeName, r.RoutingKey, false, false, amqp.Publishing{
		ContentType: contentType,
		Body:        payload,
	})

	return err
}

// func (r *RabbitMQProvider) Consume() (<-chan amqp.Delivery, error) {
// 	// for {

// 	// }
// 	// delivery, err := r.Channel.Consume(r.QueueName, r.Name, false, false, false, false, nil)

// 	// return delivery, err
// }

func (r *RabbitMQProvider) Close() error {
	if err := r.Channel.Cancel(r.Name, true); err != nil {
		return fmt.Errorf("consumer cancel failed: %s", err)
	}

	if err := r.Conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer logrus.Info("AMQP shutdown OK")

	return <-r.Done
}
