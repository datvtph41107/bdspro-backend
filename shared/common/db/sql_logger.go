package _db

import (
	"common/logging"
	"common/requestlog"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type sqlLogger struct {
	level  logger.LogLevel
	config logging.SQLConfig
}

var sqlRecorderOnce sync.Once

// NewSQLLogger dùng cùng logger process. Query mặc định chỉ được ghi khi lỗi/chậm;
// config sql.enabled bật query thường ở DEBUG. Query luôn parameterized.
// Adapter mở pool riêng của service cũng dùng hàm này để không có chính sách SQL song song.
func NewSQLLogger() logger.Interface {
	configureSQLRecorder()

	cfg := logging.SQLSettings()
	return &sqlLogger{level: logger.Info, config: cfg}
}
func (l *sqlLogger) LogMode(level logger.LogLevel) logger.Interface {
	copy := *l
	copy.level = level
	return &copy
}
func (l *sqlLogger) Info(ctx context.Context, msg string, args ...any) {
	if l.level >= logger.Info {
		requestlog.FromContext(ctx, nil).DebugContext(ctx, msg, "category", "database")
	}
}
func (l *sqlLogger) Warn(ctx context.Context, msg string, args ...any) {
	if l.level >= logger.Warn {
		requestlog.FromContext(ctx, nil).WarnContext(ctx, msg, "category", "database")
	}
}
func (l *sqlLogger) Error(ctx context.Context, msg string, args ...any) {
	if l.level >= logger.Error {
		requestlog.FromContext(ctx, nil).ErrorContext(ctx, msg, "category", "database")
	}
}

// ParamsFilter được GORM gọi trước khi dựng SQL log; giữ placeholder, không đưa giá trị cá nhân vào file.
func (l *sqlLogger) ParamsFilter(_ context.Context, sql string, _ ...any) (string, []any) {
	return sql, nil
}
func (l *sqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	// GORM LogMode: Silent < Error < Warn < Info. Higher = more verbose.
	if l.level <= logger.Silent {
		return
	}
	elapsed := time.Since(begin)
	isError := err != nil && !errors.Is(err, gorm.ErrRecordNotFound)
	isSlow := elapsed > l.config.SlowThreshold

	var level slog.Level
	switch {
	case isError:
		level = slog.LevelError
	case isSlow:
		if l.level < logger.Warn {
			return
		}
		level = slog.LevelWarn
	default:
		// Successful non-slow queries only when Info mode and sql.enabled.
		if l.level < logger.Info || !l.config.Enabled {
			return
		}
		level = slog.LevelDebug
	}
	sql, rows := fc()
	operation := "unknown"
	if parts := strings.Fields(sql); len(parts) > 0 {
		operation = strings.ToUpper(parts[0])
	}
	attrs := []any{"category", "database", "event_name", "database.query.completed", "db_system", "postgresql", "sql_operation", operation, "sql", sql, "rows", rows, "duration_ms", float64(elapsed.Microseconds()) / 1000, "query_source", sqlSource()}
	if err != nil {
		attrs = append(attrs, "error", err.Error())
	}
	requestlog.FromContext(ctx, nil).Log(ctx, level, "SQL query completed", attrs...)
}

// Bỏ các frame logger/GORM để source trỏ tới câu truy vấn ở repository.
func sqlSource() string {
	pcs := make([]uintptr, 24)
	n := runtime.Callers(2, pcs)
	frames := runtime.CallersFrames(pcs[:n])
	for {
		frame, more := frames.Next()
		if !strings.Contains(frame.File, "/gorm.io/") && !strings.Contains(frame.File, "/common/db/sql_logger.go") && !strings.HasPrefix(frame.Function, "runtime.") {
			return fmt.Sprintf("%s:%d", frame.File, frame.Line)
		}
		if !more {
			break
		}
	}
	return ""
}

// GORM Scan dùng traceRecorder tạm thay logger của DB; ParamsFilter trên sqlLogger
// không được gọi ở đường này. Cấu hình hook chính thức một lần trước khi pool nhận request.
func configureSQLRecorder() {
	sqlRecorderOnce.Do(func() {
		logger.RecorderParamsFilter = func(_ context.Context, sql string, _ ...interface{}) (string, []interface{}) {
			return sql, nil
		}
	})
}
