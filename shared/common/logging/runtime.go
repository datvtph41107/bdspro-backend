package logging

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

var currentConfig atomic.Pointer[Config]

// MustStart được gọi một lần tại main; cấu hình sai phải dừng trước khi nhận request.
// Không có config file thì chỉ xuất stdout, để container runtime sở hữu retention.
func MustStart(service string) func() {
	c, err := Load(os.Getenv("QHPRO_LOG_CONFIG"), service)
	if err != nil {
		fmt.Fprintln(os.Stderr, "logging configuration:", err)
		os.Exit(1)
	}
	closeLogs, err := Start(service, c)
	if err != nil {
		fmt.Fprintln(os.Stderr, "logging startup:", err)
		os.Exit(1)
	}
	return func() {
		if err := closeLogs(); err != nil {
			fmt.Fprintln(os.Stderr, "logging shutdown:", err)
		}
	}
}

func Start(service string, c Config) (func() error, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	if service == "" || filepath.Base(service) != service || service == "." || service == ".." {
		return nil, fmt.Errorf("invalid logging service %q", service)
	}
	runID := fmt.Sprintf("%s-%d", time.Now().UTC().Format("20060102T150405.000000000Z"), os.Getpid())
	directory := filepath.Join(c.Root, service)
	if c.Mode == "session" {
		directory = filepath.Join(directory, "runs", runID)
	}
	var release func() error = func() error { return nil }
	if c.Files {
		var err error
		release, err = lockDirectory(directory)
		if err != nil {
			return nil, err
		}
		if c.Mode == "truncate" {
			if err := truncateActiveFiles(directory); err != nil {
				_ = release()
				return nil, err
			}
		}
		if c.Mode == "session" {
			pruneRuns(filepath.Join(c.Root, service, "runs"), directory, c.KeepRuns, c.MaxAgeDays)
		}
	}
	h := newHandler(c, os.Stdout, directory)
	if h.sinks != nil {
		if _, err := h.sinks.get("info", "lifecycle"); err != nil {
			_ = release()
			return nil, err
		}
	}
	logger := slog.New(h).With("service_name", service+"-service", "environment", os.Getenv("QHPRO_ENVIRONMENT"), "run_id", runID)
	slog.SetDefault(logger)
	currentConfig.Store(&c)
	logger.Info("Process logging ready", "category", "lifecycle", "event_name", "process.logging.ready", "log_directory", directory, "mode", c.Mode, "files_enabled", c.Files, "minimum_level", c.Level)
	return func() error {
		logger.Info("Process logging stopped", "category", "lifecycle", "event_name", "process.logging.stopped")
		if h.sinks != nil {
			err := h.sinks.close()
			_ = release()
			return err
		}
		return release()
	}, nil
}

// SQLSettings trả bản sao cấu hình cho adapter GORM; core logging không import ORM.
func SQLSettings() SQLConfig {
	if c := currentConfig.Load(); c != nil {
		return c.SQL
	}
	return defaults().SQL
}
