package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

type queryLogger struct {
	logger *zap.SugaredLogger
}

func (ql *queryLogger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	ql.logger.Debugw(msg, "level", level, zap.Any("data", data))
}
