package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	var err error
	cfg, err := LoadDatabaseCredentialsFromEnv()

	if err != nil {
		log.Fatal(err)
	}

	worker, err := NewWorker(cfg)

	defer worker.Close()

	worker.Consume()
}

func (w *Worker) handle(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		log.Printf(
			"got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)
		d.Ack(false)
	}
	log.Printf("handle: deliveries channel closed")
	done <- nil
}
