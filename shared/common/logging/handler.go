package logging

import (
	"context"
	"log/slog"
	"strings"
)

type multiHandler struct {
	handlers []slog.Handler
}

func (h multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h multiHandler) Handle(ctx context.Context, record slog.Record) error {
	var first error
	for _, handler := range h.handlers {
		if !handler.Enabled(ctx, record.Level) {
			continue
		}
		if err := handler.Handle(ctx, record.Clone()); err != nil && first == nil {
			first = err
		}
	}
	return first
}

func (h multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	next := make([]slog.Handler, 0, len(h.handlers))
	for _, handler := range h.handlers {
		next = append(next, handler.WithAttrs(attrs))
	}
	return multiHandler{handlers: next}
}

func (h multiHandler) WithGroup(name string) slog.Handler {
	next := make([]slog.Handler, 0, len(h.handlers))
	for _, handler := range h.handlers {
		next = append(next, handler.WithGroup(name))
	}
	return multiHandler{handlers: next}
}

type filterHandler struct {
	handler slog.Handler
	accept  func(slog.Record) bool
}

func (h filterHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h filterHandler) Handle(ctx context.Context, record slog.Record) error {
	if h.accept != nil && !h.accept(record) {
		return nil
	}
	return h.handler.Handle(ctx, record)
}

func (h filterHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return filterHandler{handler: h.handler.WithAttrs(attrs), accept: h.accept}
}

func (h filterHandler) WithGroup(name string) slog.Handler {
	return filterHandler{handler: h.handler.WithGroup(name), accept: h.accept}
}

func channelFilter(channel string) func(slog.Record) bool {
	return func(record slog.Record) bool {
		matched := false
		record.Attrs(func(attr slog.Attr) bool {
			if attr.Key == "channel" && attr.Value.Kind() == slog.KindString && attr.Value.String() == channel {
				matched = true
				return false
			}
			return true
		})
		return matched
	}
}

func errorFilter(record slog.Record) bool {
	return record.Level >= slog.LevelError
}

func redactAttr(_ []string, attr slog.Attr) slog.Attr {
	key := normalizeSensitiveKey(attr.Key)
	for _, fragment := range []string{
		"authorization", "credential", "password", "passwd", "secret", "token", "cookie", "api_key", "apikey",
	} {
		if strings.Contains(key, fragment) {
			return slog.String(attr.Key, "[REDACTED]")
		}
	}
	return attr
}

func normalizeSensitiveKey(key string) string {
	key = strings.ToLower(strings.TrimSpace(key))
	key = strings.ReplaceAll(key, "-", "_")
	key = strings.ReplaceAll(key, ".", "_")
	return key
}
