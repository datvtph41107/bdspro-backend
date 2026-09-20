package logging

import (
	"context"
	"log"
	"log/slog"
	"strings"
)

type slogWriter struct {
	logger *slog.Logger
}

func (w slogWriter) Write(p []byte) (int, error) {
	message := strings.TrimSpace(string(p))
	if message != "" {
		w.logger.Info(message)
	}
	return len(p), nil
}

// StdLogger bridges dependencies that require *log.Logger into the canonical
// process-wide slog pipeline. Application code should continue to use log/slog.
func StdLogger(component string) *log.Logger {
	return log.New(
		slogWriter{logger: WithComponent(context.Background(), component)},
		"",
		0,
	)
}
