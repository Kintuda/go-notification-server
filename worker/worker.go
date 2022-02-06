package worker

import (
	"encoding/json"
	"strings"
	"time"

	config "github.com/Kintuda/notification-server/config"
	db "github.com/Kintuda/notification-server/database"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Worker struct {
	Name      string
	QueueName string
	Logger    *zap.Logger
	Conn      *amqp.Connection
	Channel   *amqp.Channel
	Done      chan error
}

func NewWorker(cfg *config.WorkerConfig) (*Worker, error) {
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

	if err != nil {
		return nil, err
	}

	w.Channel, err = w.Conn.Channel()

	if err != nil {
		return nil, err
	}

	if cfg.CreateExchange {
		w.Logger.Info("creating exchange", zap.String("exchange", cfg.ExchangeName))
		err = w.Channel.ExchangeDeclare(cfg.ExchangeName, "direct", true, false, false, false, nil)
	}

	if cfg.CreateQueue {
		w.Logger.Info("creating queue", zap.String("queue", cfg.QueueName))
		_, err = w.Channel.QueueDeclare(cfg.QueueName, true, false, false, false, nil)
	}

	if err != nil {
		return nil, err
	}

	routingKeys := strings.Split(cfg.RoutingKeys, ",")

	for _, routingKey := range routingKeys {
		w.Logger.Info("binding queue to consumer routing key", zap.String("routing", routingKey))
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
	if err := w.Channel.Cancel(w.Name, true); err != nil {
		return err
	}

	if err := w.Conn.Close(); err != nil {
		return err
	}

	return <-w.Done
}

func (w *Worker) ParseEvent(msg amqp.Delivery, conn *db.DatabaseConnection) {
	startTime := time.Now()

	var event Event
	err := json.Unmarshal(msg.Body, &event)

	if err != nil {
		w.Logger.Error("unmarshalling body", zap.String("body", string(msg.Body)))
		msg.Nack(true, false)
		return
	}

	if event.Id == "" {
		w.Logger.Info("empty event", zap.String("id", event.Id))
		msg.Nack(true, false)
		return
	}

	delivered := CreateWebhookTransaction(event, conn)
	w.Logger.Info("processing finished", zap.Int64("ms", time.Since(startTime).Microseconds()), zap.Bool("delivered", delivered))

	if delivered {
		msg.Ack(false)
		return
	}

	msg.Nack(false, true)
}
