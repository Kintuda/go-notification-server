package cmd

import (
	"github.com/Kintuda/notification-server/pkg/config"
	"github.com/Kintuda/notification-server/pkg/http"
	"github.com/Kintuda/notification-server/pkg/postgres"
	"github.com/Kintuda/notification-server/pkg/queue"
	"github.com/spf13/cobra"
)

func NewServerCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "server",
		Short: "start appname API",
	}

	listen := &cobra.Command{
		Use:  "listen",
		RunE: startCmd,
	}

	rootCmd.AddCommand(listen)

	return rootCmd
}

func startCmd(cmd *cobra.Command, arg []string) error {
	cfg := &config.AppConfig{}

	if err := config.LoadConfigFromEnv(cfg); err != nil {
		return err
	}

	db, err := postgres.NewDatabaseConnection(cfg.PostgresDns)

	if err != nil {
		return err
	}

	params := queue.RabbitMQParams{
		AmqpDns:        cfg.AmqpDns,
		Name:           "notification_server",
		QueueName:      cfg.QueueName,
		ExchangeName:   cfg.ExchangeName,
		RoutingKey:     cfg.RoutingKey,
		CreateExchange: cfg.CreateExchange,
		CreateQueue:    cfg.CreateQueue,
	}

	provider, err := queue.NewRabbitMQ(params)

	if err != nil {
		return err
	}

	router, err := http.NewRouter(cfg, db, *provider)

	if err != nil {
		return err
	}

	http.RegisterRoutes(router)

	server := http.NewServer(router, cfg.HttpPort)

	err = server.Init()
	return err
}
