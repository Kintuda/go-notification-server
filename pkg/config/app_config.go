package config

type AppConfig struct {
	Env            string `mapstructure:"ENV"`
	HttpPort       string `mapstructure:"HTTP_PORT"`
	PostgresDns    string `mapstructure:"POSTGRES_DNS"`
	AmqpDns        string `mapstructure:"AMQP_DNS"`
	QueueName      string `mapstructure:"AMQP_QUEUE_NAME"`
	ExchangeName   string `mapstructure:"AMQP_EXCHANGE_NAME"`
	RoutingKey     string `mapstructure:"AMQP_ROUTING_KEY"`
	CreateExchange bool   `mapstructure:"AMQP_CREATE_EXCHANGE"`
	CreateQueue    bool   `mapstructure:"AMQP_CREATE_QUEUE"`
}
