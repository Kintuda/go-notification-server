package database

import (
	"context"

	"github.com/Kintuda/notification-server/config"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type DatabaseConnection struct {
	Conn *pgxpool.Pool
}

func NewDatabaseConnection(c *config.WorkerConfig) (*DatabaseConnection, error) {
	var err error

	config, err := pgxpool.ParseConfig(c.DatabaseUrl)

	if err != nil {
		return nil, err
	}

	if c.DatabaseDebug {
		logger, err := zap.NewProduction()

		if err != nil {
			return nil, err
		}

		config.ConnConfig.Logger = zapadapter.NewLogger(logger)
	}

	conn, err := pgxpool.ConnectConfig(context.Background(), config)

	if err != nil {
		return nil, err
	}

	return &DatabaseConnection{Conn: conn}, nil
}

func (d *DatabaseConnection) CloseConnection() {
	d.Conn.Close()
}
