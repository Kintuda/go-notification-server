package main

import (
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Worker struct {
	Name      string
	QueueName string
	Logger    *zap.Logger
	Conn      *amqp.Connection
	Channel   *amqp.Channel
}

func NewWorker(cfg *WorkerConfig) (*Worker, error) {
	var err error
	l, err := zap.NewProduction()

	if err != nil {
		return nil, err
	}

	w := &Worker{
		Conn:      nil,
		Channel:   nil,
		Name:      cfg.Name,
		QueueName: cfg.QueueName,
		Logger:    l,
	}

	w.Conn, err = amqp.Dial(cfg.AmqpDns)

	defer w.Conn.Close()

	if err != nil {
		return nil, err
	}

	w.Channel, err = w.Conn.Channel()

	if err != nil {
		return nil, err
	}

	routingKeys := strings.Split(cfg.RoutingKeys, ",")

	for _, routingKey := range routingKeys {
		err := w.Channel.QueueBind(
			cfg.QueueName,
			routingKey,
			cfg.ExchangeName,
			false,
			nil,
		)

		if err != nil {
			return nil, err
		}
	}

	return w, nil
}

func (w *Worker) Consume() (<-chan amqp.Delivery, error) {
	delivery, err := w.Channel.Consume(w.QueueName, w.Name, false, false, false, false, nil)

	return delivery, err
}

func (w *Worker) Close() error {
	var err error
	err = w.Channel.Close()

	if err != nil {
		return err
	}

	err = w.Conn.Close()

	if err != nil {
		return err
	}

	return nil
}
