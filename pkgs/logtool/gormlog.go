package logtool

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

type loggerAdapter struct {
	logger *zap.Logger
}

func (a *loggerAdapter) LogMode(level logger.LogLevel) logger.Interface {
	return a
}

func (a *loggerAdapter) Info(ctx context.Context, msg string, data ...interface{}) {
	a.logger.Info("", zap.Any(msg, data))
}

func (a *loggerAdapter) Warn(ctx context.Context, msg string, data ...interface{}) {
	a.logger.Warn("", zap.Any(msg, data))
}

func (a *loggerAdapter) Error(ctx context.Context, msg string, data ...interface{}) {
	a.logger.Error("", zap.Any(msg, data))
}

func (a *loggerAdapter) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	a.logger.Debug("gorm trace", zap.String("sql", sql), zap.Int64("rows", rows), zap.Duration("elapsed", elapsed), zap.Error(err))
}
