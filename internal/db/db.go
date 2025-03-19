package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
    "github.com/amarantec/move-easy/pkg/logger"
)

var Conn *pgxpool.Pool

func OpenConnection(ctx context.Context, connectionString string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

    cfg.ConnConfig.Tracer = &tracelog.TraceLog{
        Logger: &logger.PgxLogger{},
        LogLevel: tracelog.LogLevelDebug,
    }

	Conn, err = pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

    createTables(ctx)

	return Conn, nil
}
