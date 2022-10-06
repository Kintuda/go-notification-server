package postgres

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type Repositories interface {
	NotificationRepositoryPg
}

type QueryExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

type PostgresRepository struct {
	Connection *Pool
	Trx        *pgx.Tx
}

func NewPostgresRepository(conn *Pool) *PostgresRepository {
	return &PostgresRepository{
		Connection: conn,
	}
}

func NewRepository[K Repositories](p *PostgresRepository) *K {
	if p.Trx == nil {
		return &K{Executor: p.Connection.Conn}
	}

	return &K{Executor: *p.Trx}
}
