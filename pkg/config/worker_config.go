package config

type WorkerConfig struct {
	Env            string `mapstructure:"ENV"`
	AmqpDns        string `mapstructure:"AMQP_URL"`
	ExchangeName   string `mapstructure:"AMQP_EXCHANGE_NAME"`
	QueueName      string `mapstructure:"AMQP_QUEUE_NAME"`
	RoutingKeys    string `mapstructure:"AMQP_ROUTING_KEYS"`
	Name           string `mapstructure:"AMQP_CONSUMER_NAME"`
	CreateExchange bool   `mapstructure:"AMQP_CREATE_EXCHANGE"`
	CreateQueue    bool   `mapstructure:"AMQP_CREATE_QUEUE"`
}
