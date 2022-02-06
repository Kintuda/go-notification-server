package database

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type DatabaseConnection struct {
	Conn *pgx.Conn
}

func NewDatabaseConnection(dns string) (*DatabaseConnection, error) {
	var err error
	database := &DatabaseConnection{}

	database.Conn, err = pgx.Connect(context.Background(), dns)

	if err != nil {
		return nil, err
	}

	return database, nil
}

func (d *DatabaseConnection) Close() error {
	return d.Conn.Close(context.Background())
}
