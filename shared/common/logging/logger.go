// Package logging is the single backend logging owner.
//
// Application code uses log/slog only. This package owns configuration,
// service identity, correlation enrichment, redaction and local file routing.
package logging

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"common/requestlog"
)

const (
	defaultNamespace = "bdspro"
	defaultLogRoot   = ".tmp/development/logs"
)

type Config struct {
	ServiceName   string
	Namespace     string
	Environment   string
	InstanceID    string
	Level         slog.Level
	Output        string
	Root          string
	FileMode      string
	ConsoleFormat string
	AddSource     bool
	Channels      []string
	RunID         string
}

func FromEnv(serviceName string) Config {
	return Config{
		ServiceName:   strings.TrimSpace(serviceName),
		Namespace:     envOr("QHPRO_LOG_SERVICE_NAMESPACE", defaultNamespace),
		Environment:   envOr("ENVIRONMENT", "development"),
		InstanceID:    envOr("QHPRO_SERVICE_INSTANCE_ID", defaultInstanceID()),
		Level:         parseLevel(os.Getenv("QHPRO_LOG_LEVEL")),
		Output:        envOr("QHPRO_LOG_OUTPUT", "both"),
		Root:          envOr("QHPRO_LOG_ROOT", defaultLogRoot),
		FileMode:      envOr("QHPRO_LOG_FILE_MODE", "session"),
		ConsoleFormat: envOr("QHPRO_LOG_CONSOLE_FORMAT", "text"),
		AddSource:     boolEnv("QHPRO_LOG_SOURCE", false),
		Channels:      splitCSV(envOr("QHPRO_LOG_CHANNELS", "http,sql,external,audit")),
		RunID:         strings.TrimSpace(os.Getenv("QHPRO_LOG_RUN_ID")),
	}
}

// Configure installs the one process-wide default logger. Composition roots call
// this once and defer the returned close function.
func Configure(serviceName string) (func() error, error) {
	logger, closeLogger, err := New(FromEnv(serviceName))
	if err != nil {
		return nil, err
	}
	slog.SetDefault(logger)
	return closeLogger, nil
}

func New(cfg Config) (*slog.Logger, func() error, error) {
	cfg.ServiceName = strings.TrimSpace(cfg.ServiceName)
	if cfg.ServiceName == "" {
		return nil, nil, errors.New("logging requires service name")
	}
	if cfg.Namespace == "" {
		cfg.Namespace = defaultNamespace
	}
	if cfg.Environment == "" {
		cfg.Environment = "development"
	}
	if cfg.InstanceID == "" {
		cfg.InstanceID = defaultInstanceID()
	}
	if cfg.Root == "" {
		cfg.Root = defaultLogRoot
	}

	options := &slog.HandlerOptions{Level: cfg.Level, AddSource: cfg.AddSource, ReplaceAttr: redactAttr}
	var handlers []slog.Handler
	var closers []io.Closer

	output := strings.ToLower(strings.TrimSpace(cfg.Output))
	if output == "" {
		output = "both"
	}
	if output == "stdout" || output == "both" {
		var console slog.Handler
		if strings.EqualFold(cfg.ConsoleFormat, "json") {
			console = slog.NewJSONHandler(os.Stdout, options)
		} else {
			console = slog.NewTextHandler(os.Stdout, options)
		}
		handlers = append(handlers, console)
	}

	if output == "file" || output == "both" {
		directory, flags, err := resolveFileTarget(cfg)
		if err != nil {
			return nil, nil, err
		}
		if err := os.MkdirAll(filepath.Join(directory, "channels"), 0o755); err != nil {
			return nil, nil, fmt.Errorf("create log directory: %w", err)
		}

		runtimeFile, err := os.OpenFile(filepath.Join(directory, "runtime.jsonl"), flags, 0o644)
		if err != nil {
			return nil, nil, fmt.Errorf("open runtime log: %w", err)
		}
		closers = append(closers, runtimeFile)
		handlers = append(handlers, slog.NewJSONHandler(runtimeFile, options))

		for _, channel := range cfg.Channels {
			channel = strings.TrimSpace(channel)
			if channel == "" || channel == "error" {
				continue
			}
			channelFile, openErr := os.OpenFile(filepath.Join(directory, "channels", channel+".jsonl"), flags, 0o644)
			if openErr != nil {
				closeAll(closers)
				return nil, nil, fmt.Errorf("open %s log: %w", channel, openErr)
			}
			closers = append(closers, channelFile)
			handlers = append(handlers, channelProjectionHandler{
				handler: slog.NewJSONHandler(channelFile, options),
				target:  channel,
			})
		}

		errorFile, err := os.OpenFile(filepath.Join(directory, "channels", "error.jsonl"), flags, 0o644)
		if err != nil {
			closeAll(closers)
			return nil, nil, fmt.Errorf("open error log: %w", err)
		}
		closers = append(closers, errorFile)
		handlers = append(handlers, filterHandler{handler: slog.NewJSONHandler(errorFile, options), accept: errorFilter})
	}

	if len(handlers) == 0 {
		return nil, nil, fmt.Errorf("unsupported QHPRO_LOG_OUTPUT %q", cfg.Output)
	}

	logger := slog.New(multiHandler{handlers: handlers}).With(
		slog.String("service.namespace", cfg.Namespace),
		slog.String("service.name", cfg.ServiceName),
		slog.String("service.instance.id", cfg.InstanceID),
		slog.String("deployment.environment", cfg.Environment),
	)
	return logger, func() error { return closeAll(closers) }, nil
}

// FromContext enriches the configured logger with the canonical request and
// operation identifiers already carried by common/request metadata.
func FromContext(ctx context.Context) *slog.Logger {
	return requestlog.FromContext(ctx, slog.Default())
}

func WithComponent(ctx context.Context, component string) *slog.Logger {
	logger := FromContext(ctx)
	if component = strings.TrimSpace(component); component != "" {
		logger = logger.With(slog.String("component", component))
	}
	return logger
}

func WithChannel(ctx context.Context, channel string) *slog.Logger {
	logger := FromContext(ctx)
	if channel = strings.TrimSpace(channel); channel != "" {
		logger = logger.With(slog.String("channel", channel))
	}
	return logger
}

func resolveFileTarget(cfg Config) (string, int, error) {
	mode := strings.ToLower(strings.TrimSpace(cfg.FileMode))
	serviceRoot := filepath.Join(cfg.Root, cfg.ServiceName)
	switch mode {
	case "append":
		return filepath.Join(serviceRoot, "current"), os.O_CREATE | os.O_WRONLY | os.O_APPEND, nil
	case "truncate":
		return filepath.Join(serviceRoot, "current"), os.O_CREATE | os.O_WRONLY | os.O_TRUNC, nil
	case "session", "":
		runID := cfg.RunID
		if runID == "" {
			runID = time.Now().UTC().Format("20060102T150405.000000000Z") + "-pid-" + strconv.Itoa(os.Getpid())
		}
		return filepath.Join(serviceRoot, "runs", sanitizeSegment(runID)), os.O_CREATE | os.O_WRONLY | os.O_APPEND, nil
	default:
		return "", 0, fmt.Errorf("unsupported QHPRO_LOG_FILE_MODE %q", cfg.FileMode)
	}
}

func parseLevel(value string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	seen := map[string]struct{}{}
	for _, part := range parts {
		part = strings.ToLower(strings.TrimSpace(part))
		if part == "" {
			continue
		}
		if _, ok := seen[part]; ok {
			continue
		}
		seen[part] = struct{}{}
		result = append(result, part)
	}
	return result
}

func sanitizeSegment(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, string(filepath.Separator), "_")
	value = strings.ReplaceAll(value, "..", "_")
	if value == "" {
		return "unknown"
	}
	return value
}

func closeAll(closers []io.Closer) error {
	var result error
	for i := len(closers) - 1; i >= 0; i-- {
		if err := closers[i].Close(); err != nil && result == nil {
			result = err
		}
	}
	return result
}

func defaultInstanceID() string {
	host, _ := os.Hostname()
	if host == "" {
		host = "unknown"
	}
	return host + ":" + strconv.Itoa(os.Getpid())
}

func envOr(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func boolEnv(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}
