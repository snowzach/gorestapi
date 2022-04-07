package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"
)

// The Logger Inferface for the database driver.
type Logger interface {
	Printf(template string, args ...interface{})
}

type nopLogger struct{}

func (_ nopLogger) Printf(template string, args ...interface{}) {}

type queryLogger struct {
	Logger
}

func (ql *queryLogger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	ql.Printf("%s: %v", msg, data)
}
