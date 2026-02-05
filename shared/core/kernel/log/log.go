package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config for creating logger (you may keep factory.go separately)
type Config struct {
	Service string
	Env     string // develop|staging|prod
	Level   string // debug|info|warn|error
	Format  string // json|console
}

const (
	FieldService = "service"
	FieldTrace   = "trace"
	FieldCode    = "code"
	FieldTag     = "tag"
)

// Logger wraps zap.SugaredLogger (or zap.Logger) with convenience helpers
type Logger struct {
	*zap.SugaredLogger
	service string
}

// New builds logger from config
func New(cfg Config) (*Logger, error) {
	level := zapcore.InfoLevel
	if l, err := zapcore.ParseLevel(cfg.Level); err == nil {
		level = l
	}
	enc := "json"
	if cfg.Format != "" {
		enc = cfg.Format
	}
	zcfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      cfg.Env != "prod",
		Encoding:         enc,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:     "ts",
			LevelKey:    "level",
			MessageKey:  "msg",
			EncodeTime:  zapcore.ISO8601TimeEncoder,
			EncodeLevel: zapcore.LowercaseLevelEncoder,
		},
	}
	z, err := zcfg.Build()
	if err != nil {
		return nil, err
	}
	return &Logger{
		SugaredLogger: z.Sugar(),
		service:       cfg.Service,
	}, nil
}

// Nop returns a no-op logger (safest fallback)
func Nop() *Logger {
	z := zap.NewNop().Sugar()
	return &Logger{SugaredLogger: z, service: "nop"}
}

// WithTag returns a logger with a tag field
func (l *Logger) WithTag(tag string) *Logger {
	return &Logger{SugaredLogger: l.With(FieldTag, tag), service: l.service}
}

// WithTrace returns a logger with trace id (for request-scoped logs)
func (l *Logger) WithTrace(trace string) *Logger {
	return &Logger{SugaredLogger: l.With(FieldTrace, trace), service: l.service}
}

// Info with code and message (structured)
func (l *Logger) Info(code string, msg string, keysAndValues ...interface{}) {
	args := append([]interface{}{FieldService, l.service, FieldCode, code}, keysAndValues...)
	l.SugaredLogger.Infow(msg, args...)
}

// Error with err and structured fields
func (l *Logger) Error(code string, msg string, err error, keysAndValues ...interface{}) {
	args := append([]interface{}{FieldService, l.service, FieldCode, code, "error", err}, keysAndValues...)
	l.SugaredLogger.Errorw(msg, args...)
}

// Debug helper
func (l *Logger) Debug(msg string, keysAndValues ...interface{}) {
	args := append([]interface{}{FieldService, l.service}, keysAndValues...)
	l.SugaredLogger.Debugw(msg, args...)
}
