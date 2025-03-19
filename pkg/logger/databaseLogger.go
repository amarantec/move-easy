package logger

import (
    "log"
    "github.com/jackc/pgx/v5/tracelog"
    "context"
)

type PgxLogger struct{}

func (l *PgxLogger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
    log.Printf("[DATABASE] [%s] %s - %v", level, msg, data)
}
