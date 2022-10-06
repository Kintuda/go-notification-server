package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type Pool struct {
	Conn *pgxpool.Pool
}

func NewDatabaseConnection(dns string) (*Pool, error) {
	config, err := pgxpool.ParseConfig(dns)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse config: %v\n", err)
		os.Exit(1)
	}

	logrusLogger := &logrus.Logger{
		Out:          os.Stderr,
		Formatter:    new(logrus.JSONFormatter),
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}

	config.ConnConfig.Logger = logrusadapter.NewLogger(logrusLogger)

	conn, err := pgxpool.ConnectConfig(context.Background(), config)

	if err != nil {
		return nil, err
	}

	return &Pool{Conn: conn}, nil
}

func (d *Pool) CloseConnection() {
	d.Conn.Close()
}
