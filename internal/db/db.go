package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
    "github.com/amarantec/move-easy/pkg/logger"
    "log"
    "time"
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

    for {
        select {
        case <-ctx.Done():
            return nil, fmt.Errorf("Timed out, trying to connect to the database: %w", ctx.Err())
        default:
	        Conn, err = pgxpool.NewWithConfig(ctx, cfg)
	        if err == nil {
                if err := Conn.Ping(ctx); err == nil {
                    createTables(ctx)
                    return Conn, nil
                }
                log.Printf("Database not yet available, trying again in 2 seconds... (%v)\n", err)
            } else {
                log.Printf("Error trying to connect to the database, trying again in 2 seconds... (%v)\n", err)
           }
           time.Sleep(2 * time.Second)

        }
    }
}
