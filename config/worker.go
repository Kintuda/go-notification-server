package config

import (
	"github.com/Netflix/go-env"
)

type WorkerConfig struct {
	Env            string `env:"ENV,required=true"`
	AmqpDns        string `env:"AMQP_URL,required=true"`
	ExchangeName   string `env:"AMQP_EXCHANGE_NAME,required=true"`
	QueueName      string `env:"AMQP_QUEUE_NAME,required=true"`
	RoutingKeys    string `env:"AMQP_ROUTING_KEYS,required=true"`
	Name           string `env:"AMQP_CONSUMER_NAME,required=true"`
	CreateExchange bool   `env:"AMQP_CREATE_EXCHANGE,required=true"`
	CreateQueue    bool   `env:"AMQP_CREATE_QUEUE,required=true"`
	DatabaseUrl    string `env:"DATABASE_URL,required=true"`
	DatabaseDebug  bool   `env:"DATABASE_DEBUG,default=false"`
}

func LoadConfigFromEnv() (*WorkerConfig, error) {
	var cfg WorkerConfig

	if _, err := env.UnmarshalFromEnviron(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
