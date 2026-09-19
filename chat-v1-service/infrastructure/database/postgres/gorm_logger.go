package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"common/logging"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type slogGormLogger struct {
	level                gormlogger.LogLevel
	slowThreshold        time.Duration
	ignoreRecordNotFound bool
}

func newGormLogger() gormlogger.Interface {
	return slogGormLogger{
		level:                gormlogger.Info,
		slowThreshold:        200 * time.Millisecond,
		ignoreRecordNotFound: true,
	}
}

func (l slogGormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	l.level = level
	return l
}

func (l slogGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.level < gormlogger.Info {
		return
	}

	logging.WithComponent(
		ctx,
		"postgres",
	).Info(
		fmt.Sprintf(msg, data...),
	)
}

func (l slogGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.level < gormlogger.Warn {
		return
	}

	logging.WithComponent(
		ctx,
		"postgres",
	).Warn(
		fmt.Sprintf(msg, data...),
	)
}

func (l slogGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.level < gormlogger.Error {
		return
	}

	logging.WithComponent(
		ctx,
		"postgres",
	).Error(
		fmt.Sprintf(msg, data...),
	)
}

func (l slogGormLogger) Trace(
	ctx context.Context,
	begin time.Time,
	fc func() (sql string, rowsAffected int64),
	err error,
) {
	if l.level == gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)

	switch {
	case err != nil &&
		l.level >= gormlogger.Error &&
		(!l.ignoreRecordNotFound ||
			!errors.Is(err, gorm.ErrRecordNotFound)):

		sql, rows := fc()

		logging.WithComponent(
			ctx,
			"postgres",
		).Error(
			"database query failed",
			slog.Duration("duration", elapsed),
			slog.Int64("rows", rows),
			slog.String("sql", sql),
			slog.Any("error", err),
		)

	case l.slowThreshold > 0 &&
		elapsed > l.slowThreshold &&
		l.level >= gormlogger.Warn:

		sql, rows := fc()

		logging.WithComponent(
			ctx,
			"postgres",
		).Warn(
			"slow database query",
			slog.Duration("duration", elapsed),
			slog.Duration(
				"slow_threshold",
				l.slowThreshold,
			),
			slog.Int64("rows", rows),
			slog.String("sql", sql),
		)

	case l.level >= gormlogger.Info:
		sql, rows := fc()

		logging.WithComponent(
			ctx,
			"postgres",
		).Info(
			"database query",
			slog.Duration("duration", elapsed),
			slog.Int64("rows", rows),
			slog.String("sql", sql),
		)
	}
}
