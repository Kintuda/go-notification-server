package config

import (
	"errors"

	"github.com/Netflix/go-env"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

var (
	errMissingConfig = errors.New("missing configuration")
)

type WorkerConfig struct {
	Env            string `env:"ENV" validate:"required"`
	AmqpDns        string `env:"AMQP_URL" validate:"required"`
	ExchangeName   string `env:"AMQP_EXCHANGE_NAME" validate:"required"`
	QueueName      string `env:"AMQP_QUEUE_NAME" validate:"required"`
	RoutingKeys    string `env:"AMQP_ROUTING_KEYS" validate:"required"`
	Name           string `env:"AMQP_CONSUMER_NAME" validate:"required"`
	CreateExchange bool   `env:"AMQP_CREATE_EXCHANGE" validate:"required"`
	CreateQueue    bool   `env:"AMQP_CREATE_QUEUE" validate:"required"`
	DatabaseUrl    string `env:"DATABASE_URL" validate:"required"`
}

func LoadWorkerCredentialsFromEnv() (*WorkerConfig, error) {
	var cfg WorkerConfig

	if err := godotenv.Load(); err != nil {
		return nil, errMissingConfig
	}

	if _, err := env.UnmarshalFromEnviron(&cfg); err != nil {
		return nil, err
	}

	validate := validator.New()

	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
