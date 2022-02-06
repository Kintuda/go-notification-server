package main

import (
	"log"

	config "github.com/Kintuda/notification-server/config"
	db "github.com/Kintuda/notification-server/database"
	worker "github.com/Kintuda/notification-server/worker"
	"go.uber.org/zap"
)

func main() {
	var err error
	logger, err := zap.NewProduction()

	if err != nil {
		log.Fatal(err)
		return
	}

	cfg, err := config.LoadWorkerCredentialsFromEnv()

	if err != nil {
		logger.Error("error while loading envs", zap.Error(err))
		return
	}

	database, err := db.NewDatabaseConnection(cfg.DatabaseUrl)

	if err != nil {
		logger.Error("error while connecting to the database", zap.Error(err))
		return
	}
	logger.Info("Database connected")

	worker, err := worker.NewWorker(cfg)

	if err != nil {
		logger.Error("error while connecting to the amqp", zap.Error(err))
		return
	}

	logger.Info("Worker connected")

	deliveries, err := worker.Consume()

	if err != nil {
		logger.Error("error while getting messages from amqp", zap.Error(err))
		return
	}

	logger.Info("processing messages")

	forever := make(chan bool)

	go func() {
		for message := range deliveries {
			// For example, show received message in a console.
			log.Printf(" > Received message: %s\n", message.Body)
			worker.ParseEvent(message, database)
		}
	}()

	<-forever

	defer func() {
		if err := database.Close(); err != nil {
			logger.Error("error while closing database connection")
		}

		if err := worker.Close(); err != nil {
			logger.Error("error while closing amqp connection")
		}

		logger.Info("exiting process")
	}()
}
