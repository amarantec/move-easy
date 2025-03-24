package logger

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/tracelog"
)

type PgxLogger struct{}

func (l *PgxLogger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	log.Printf("[QUERY] [%s] %s - %v", level, msg, data)
}
