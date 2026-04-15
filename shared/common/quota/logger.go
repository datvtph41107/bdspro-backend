package quota

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Log levels
const (
	LevelInfo  = "INFO"
	LevelWarn  = "WARN"
	LevelError = "ERROR"
)

var (
	quotaInfoLog  *log.Logger
	quotaWarnLog  *log.Logger
	quotaErrorLog *log.Logger
)

// InitQuotaLogger khởi tạo logger cho quota (ghi file + console)
func InitQuotaLogger(logDir string, consoleOutput bool) error {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	infoFile, err := os.OpenFile(filepath.Join(logDir, "quota_info.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	warnFile, err := os.OpenFile(filepath.Join(logDir, "quota_warn.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	errorFile, err := os.OpenFile(filepath.Join(logDir, "quota_error.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	if consoleOutput {
		// Ghi đồng thời ra console và file
		quotaInfoLog = log.New(io.MultiWriter(infoFile, os.Stdout), "", 0)
		quotaWarnLog = log.New(io.MultiWriter(warnFile, os.Stdout), "", 0)
		quotaErrorLog = log.New(io.MultiWriter(errorFile, os.Stdout), "", 0)
	} else {
		quotaInfoLog = log.New(infoFile, "", 0)
		quotaWarnLog = log.New(warnFile, "", 0)
		quotaErrorLog = log.New(errorFile, "", 0)
	}
	return nil
}

func timestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// LogCheck ghi log kiểm tra quota
func LogCheck(ctx context.Context, userID, service, action, pool string, used, limit, remaining int64, resetAt time.Time, latencyMs int64) {
	if quotaInfoLog == nil {
		return
	}
	msg := fmt.Sprintf("[%s] %-5s [%s] %s.%s | pool=%s | used=%d/%d | remaining=%d | reset=%s | latency=%dms",
		timestamp(), LevelInfo, userID, service, action, pool, used, limit, remaining,
		resetAt.Format("2006-01-02"), latencyMs)
	quotaInfoLog.Println(msg)
}

// LogConsume ghi log tiêu thụ quota
func LogConsume(ctx context.Context, userID, service, action, pool string, cost, usedBefore, usedAfter, limit, remaining int64, resetAt time.Time, latencyMs int64) {
	if quotaInfoLog == nil {
		return
	}
	msg := fmt.Sprintf("[%s] %-5s [%s] %s.%s | pool=%s | cost=%d | usage: %d→%d (limit=%d) | remaining=%d | latency=%dms",
		timestamp(), LevelInfo, userID, service, action, pool, cost, usedBefore, usedAfter, limit, remaining, latencyMs)
	quotaInfoLog.Println(msg)
}

// LogReject ghi log từ chối do hết quota
func LogReject(ctx context.Context, userID, service, action, pool string, used, limit int64, reason string) {
	if quotaWarnLog == nil {
		return
	}
	msg := fmt.Sprintf("[%s] %-5s [%s] %s.%s | pool=%s | REJECTED: used=%d/%d | reason=%s",
		timestamp(), LevelWarn, userID, service, action, pool, used, limit, reason)
	quotaWarnLog.Println(msg)
}

// LogWarn ghi log cảnh báo (sắp hết quota)
func LogWarn(ctx context.Context, userID, service, action, pool string, remaining, limit int64, message string) {
	if quotaWarnLog == nil {
		return
	}
	percent := float64(remaining) / float64(limit) * 100
	msg := fmt.Sprintf("[%s] %-5s [%s] %s.%s | pool=%s | WARNING: only %d/%d (%.1f%%) remaining | %s",
		timestamp(), LevelWarn, userID, service, action, pool, remaining, limit, percent, message)
	quotaWarnLog.Println(msg)
}

// LogError ghi log lỗi kỹ thuật (Redis, policy...)
func LogError(ctx context.Context, userID, service, action string, err error, extra map[string]interface{}) {
	if quotaErrorLog == nil {
		return
	}
	extraStr := ""
	if len(extra) > 0 {
		extraStr = fmt.Sprintf(" | extra=%v", extra)
	}
	msg := fmt.Sprintf("[%s] %-5s [%s] %s.%s | ERROR: %v%s",
		timestamp(), LevelError, userID, service, action, err, extraStr)
	quotaErrorLog.Println(msg)
}
