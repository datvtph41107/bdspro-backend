package logging

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFromEnvUsesCanonicalRuntimeEnvironment(t *testing.T) {
	t.Setenv("QHPRO_ENVIRONMENT", "staging")
	t.Setenv("ENVIRONMENT", "production")

	cfg := FromEnv("gateway-service")
	if cfg.Environment != "staging" {
		t.Fatalf("environment = %q, want canonical QHPRO_ENVIRONMENT value", cfg.Environment)
	}
}

func TestFileRoutingUsesOneRecordSchemaWithChannelProjections(t *testing.T) {
	root := t.TempDir()
	logger, closeLogger, err := New(Config{
		ServiceName: "user-service",
		Namespace:   "bdspro",
		Environment: "test",
		InstanceID:  "instance-1",
		Level:       slog.LevelDebug,
		Output:      "file",
		Root:        root,
		FileMode:    "session",
		Channels:    []string{"http", "sql"},
		RunID:       "run-1",
	})
	if err != nil {
		t.Fatal(err)
	}

	logger.With(
		slog.String("channel", "http"),
		slog.String("request_id", "req-123"),
	).Info("http request completed",
		slog.String("password", "must-not-leak"),
	)
	logger.With(slog.String("channel", "sql")).Debug("sql query completed")
	logger.Error("runtime failure")

	if err := closeLogger(); err != nil {
		t.Fatal(err)
	}

	directory := filepath.Join(root, "user-service", "runs", "run-1")
	runtime := readLogFile(t, filepath.Join(directory, "runtime.jsonl"))
	httpLog := readLogFile(t, filepath.Join(directory, "channels", "http.jsonl"))
	sqlLog := readLogFile(t, filepath.Join(directory, "channels", "sql.jsonl"))
	errorLog := readLogFile(t, filepath.Join(directory, "channels", "error.jsonl"))

	for _, required := range []string{
		`"service.name":"user-service"`,
		`"service.namespace":"bdspro"`,
		`"service.instance.id":"instance-1"`,
		`"deployment.environment.name":"test"`,
		`"msg":"http request completed"`,
		`"request_id":"req-123"`,
		`"password":"[REDACTED]"`,
	} {
		if !strings.Contains(runtime, required) {
			t.Fatalf("runtime log missing %q:\n%s", required, runtime)
		}
	}
	if strings.Contains(runtime, `"deployment.environment":`) {
		t.Fatalf("runtime log contains deprecated deployment.environment key: %s", runtime)
	}
	if strings.Contains(runtime, "must-not-leak") {
		t.Fatalf("runtime log leaked secret: %s", runtime)
	}

	if !strings.Contains(httpLog, `"msg":"http request completed"`) {
		t.Fatalf("http projection missing bound-channel record: %s", httpLog)
	}
	if strings.Contains(httpLog, `"msg":"sql query completed"`) {
		t.Fatalf("http projection contains sql record: %s", httpLog)
	}
	if !strings.Contains(sqlLog, `"msg":"sql query completed"`) {
		t.Fatalf("sql projection missing sql record: %s", sqlLog)
	}
	if strings.Contains(sqlLog, `"msg":"http request completed"`) {
		t.Fatalf("sql projection contains http record: %s", sqlLog)
	}
	if !strings.Contains(errorLog, `"msg":"runtime failure"`) {
		t.Fatalf("error projection missing error-level record: %s", errorLog)
	}
}

func TestUnsupportedOutputFailsClosed(t *testing.T) {
	_, closeLogger, err := New(Config{ServiceName: "user-service", Output: "mystery"})
	if err == nil {
		if closeLogger != nil {
			_ = closeLogger()
		}
		t.Fatal("expected unsupported output to fail")
	}
}

func readLogFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}
