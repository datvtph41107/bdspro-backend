package logging

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/natefinch/lumberjack.v2"
)

type fileSink struct {
	writer  *lumberjack.Logger
	handler slog.Handler
}
type sinkSet struct {
	mu        sync.Mutex
	files     map[string]fileSink
	directory string
	config    Config
	closed    bool
}

type handlerStep struct {
	group string
	attrs []slog.Attr
}

type handler struct {
	console  slog.Handler
	sinks    *sinkSet
	config   Config
	steps    []handlerStep
	groups   []string
	category string
}

// File phân theo level/category giới hạn, không tạo file theo user/request/URL để tránh nổ số file.
func (s *sinkSet) get(level, category string) (slog.Handler, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil, errors.New("log sinks are closed")
	}
	key := level + "/" + category
	if f, ok := s.files[key]; ok {
		return f.handler, nil
	}
	path := filepath.Join(s.directory, level, category+".jsonl")
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, err
	}

	// Tạo trước để permission/error cấu hình không bị giấu tới lần Write đầu tiên.
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}
	if err = f.Close(); err != nil {
		return nil, err
	}
	writer := &lumberjack.Logger{Filename: path, MaxSize: s.config.MaxSizeMB, MaxBackups: s.config.MaxBackups, MaxAge: s.config.MaxAgeDays}
	h := slog.NewJSONHandler(writer, options(s.config))
	s.files[key] = fileSink{writer: writer, handler: h}
	return h, nil
}
func (s *sinkSet) close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
	var errs []error
	for _, f := range s.files {
		errs = append(errs, f.writer.Close())
	}
	return errors.Join(errs...)
}
func options(c Config) *slog.HandlerOptions {
	return &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: c.AddSource, ReplaceAttr: redact}
}
func redact(_ []string, a slog.Attr) slog.Attr {
	key := strings.ToLower(a.Key)
	for _, s := range []string{"password", "authorization", "cookie", "secret", "token", "api_key", "dsn"} {
		if strings.Contains(key, s) {
			return slog.String(a.Key, "[REDACTED]")
		}
	}
	if a.Value.Kind() == slog.KindString {
		a.Value = slog.StringValue(strings.ReplaceAll(a.Value.String(), "\x1b", ""))
	}
	return a
}
func (h *handler) Enabled(ctx context.Context, l slog.Level) bool {
	// Category có thể được gắn ở record, nên chỉ lọc tại Handle; không chặn mất category override.
	min, _ := parseLevel(h.config.Level)
	for _, v := range h.config.Categories {
		if v != "off" {
			n, _ := parseLevel(v)
			if n < min {
				min = n
			}
		}
	}
	return l >= min
}
func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	category := h.category
	hasCategory := category != ""
	hasEvent := false
	for _, step := range h.steps {
		for _, a := range step.attrs {
			if a.Key == "event_name" {
				hasEvent = true
			}
		}
	}
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == "event_name" {
			hasEvent = true
		}
		if a.Key == "category" && len(h.groups) == 0 {
			category = a.Value.String()
			hasCategory = true
		}
		return true
	})
	if !validCategory(category) {
		category = "application"
	}
	threshold := h.config.Level
	if v, ok := h.config.Categories[category]; ok {
		threshold = v
	}
	if threshold == "off" {
		return nil
	}
	min, _ := parseLevel(threshold)
	if r.Level < min {
		return nil
	}
	apply := func(base slog.Handler) slog.Handler {
		// Metadata mặc định ở root, không tạo group rỗng giả cho record không có attrs.
		if !hasCategory {
			base = base.WithAttrs([]slog.Attr{slog.String("category", category)})
		}
		if !hasEvent {
			base = base.WithAttrs([]slog.Attr{slog.String("event_name", category+".log")})
		}
		for _, step := range h.steps {
			if step.group != "" {
				base = base.WithGroup(step.group)
			} else {
				base = base.WithAttrs(step.attrs)
			}
		}
		return base
	}
	var errs []error
	if h.console != nil {
		errs = append(errs, apply(h.console).Handle(ctx, r))
	}
	if h.sinks != nil {
		level := "info"
		if r.Level >= slog.LevelError {
			level = "error"
		} else if r.Level >= slog.LevelWarn {
			level = "warn"
		} else if r.Level < slog.LevelInfo {
			level = "debug"
		}
		out, err := h.sinks.get(level, category)
		if err == nil {
			err = apply(out).Handle(ctx, r)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "logging write failed:", err)
		}
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}
func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	next := *h
	attrs = append([]slog.Attr(nil), attrs...)
	for i, a := range attrs {
		attrs[i].Value = a.Value.Resolve()
		if a.Key == "category" && len(h.groups) == 0 {
			next.category = a.Value.String()
		}
	}
	next.steps = append(append([]handlerStep(nil), h.steps...), handlerStep{attrs: attrs})
	return &next
}
func (h *handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	next := *h
	next.groups = append(append([]string(nil), h.groups...), name)
	next.steps = append(append([]handlerStep(nil), h.steps...), handlerStep{group: name})
	return &next
}
func newHandler(c Config, out io.Writer, directory string) *handler {
	h := &handler{config: c}
	if c.Console {
		if c.Format == "text" {
			h.console = slog.NewTextHandler(out, options(c))
		} else {
			h.console = slog.NewJSONHandler(out, options(c))
		}
	}
	if c.Files {
		h.sinks = &sinkSet{files: map[string]fileSink{}, config: c, directory: directory}
	}
	return h
}
