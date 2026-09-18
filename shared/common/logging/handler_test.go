package logging

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"testing/slogtest"
)

func TestHandlerSlogContract(t *testing.T) {
	var output bytes.Buffer
	c := defaults()
	c.Level = "debug"
	h := newHandler(c, &output, "")
	results := func() []map[string]any {
		var out []map[string]any
		for _, line := range bytes.Split(bytes.TrimSpace(output.Bytes()), []byte("\n")) {
			var event map[string]any
			if err := json.Unmarshal(line, &event); err != nil {
				t.Fatal(err)
			}
			out = append(out, event)
		}
		return out
	}
	if err := slogtest.TestHandler(h, results); err != nil {
		t.Fatal(err)
	}
}

func TestFileRoutingConcurrentAndRedaction(t *testing.T) {
	c := defaults()
	c.Console = false
	c.Files = true
	c.Root = t.TempDir()
	c.Level = "debug"
	h := newHandler(c, nil, c.Root)
	log := slog.New(h).With("category", "database")
	var wg sync.WaitGroup
	for i := 0; i < 40; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Error("query failed", "password", "secret-value", "sql_operation", "SELECT")
		}()
	}
	wg.Wait()
	if err := h.sinks.close(); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(c.Root, "error", "database.jsonl")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Count(data, []byte("\n")) != 40 || bytes.Contains(data, []byte("secret-value")) {
		t.Fatalf("invalid output: %s", data)
	}
	if _, err = os.Stat(filepath.Join(c.Root, "info", "database.jsonl")); !os.IsNotExist(err) {
		t.Fatal("error must not duplicate into info")
	}
}

func TestConfigOverridesAndTypo(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "logging.yml")
	raw := `version: 1
defaults:
  files: true
  root: logs
  categories:
    grpc: warn
    database: info
services:
  tqd:
    categories:
      database: debug
    sql:
      enabled: true
`
	os.WriteFile(path, []byte(raw), 0600)
	c, err := Load(path, "tqd")
	if err != nil {
		t.Fatal(err)
	}
	if c.Categories["grpc"] != "warn" || c.Categories["database"] != "debug" || !c.SQL.Enabled || c.Root != filepath.Join(dir, "logs") {
		t.Fatalf("config=%+v", c)
	}
	os.WriteFile(path, []byte(strings.Replace(raw, "files:", "filess:", 1)), 0600)
	if _, err := Load(path, "tqd"); err == nil {
		t.Fatal("unknown field accepted")
	}
}

func TestRotationTruncateAndWriterLock(t *testing.T) {
	c := defaults()
	c.Console = false
	c.Files = true
	c.Root = t.TempDir()
	c.MaxSizeMB = 1
	release, err := lockDirectory(c.Root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := lockDirectory(c.Root); err == nil {
		t.Fatal("second writer accepted")
	}
	release()
	h := newHandler(c, nil, c.Root)
	log := slog.New(h)
	log.Info(strings.Repeat("a", 700000))
	log.Info(strings.Repeat("b", 700000))
	h.sinks.close()
	files, _ := filepath.Glob(filepath.Join(c.Root, "info", "application*.jsonl"))
	if len(files) < 2 {
		t.Fatalf("rotation files=%v", files)
	}
	c.Mode = "truncate"
	if err := truncateActiveFiles(c.Root); err != nil {
		t.Fatal(err)
	}
	h = newHandler(c, nil, c.Root)
	slog.New(h).Info("new run")
	h.sinks.close()
	data, _ := os.ReadFile(filepath.Join(c.Root, "info", "application.jsonl"))
	if bytes.Contains(data, []byte("aaaa")) || !bytes.Contains(data, []byte("new run")) {
		t.Fatal("truncate did not replace active file")
	}
}

func TestRunRetentionPreservesActiveAndUnownedDirectories(t *testing.T) {
	root := t.TempDir()
	names := []string{"20260914T000000.000000000Z-1", "20260914T000001.000000000Z-2", "20260914T000002.000000000Z-3"}
	for _, name := range names {
		os.MkdirAll(filepath.Join(root, name), 0700)
	}
	active := filepath.Join(root, names[0])
	release, err := lockDirectory(active)
	if err != nil {
		t.Fatal(err)
	}
	defer release()
	unrelated := filepath.Join(root, "operator-notes")
	os.MkdirAll(unrelated, 0700)
	current := filepath.Join(root, names[2])
	pruneRuns(root, current, 1, 7)
	for _, path := range []string{active, current, unrelated} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("must keep %s: %v", path, err)
		}
	}
	if _, err := os.Stat(filepath.Join(root, names[1])); !os.IsNotExist(err) {
		t.Fatal("old inactive run not removed")
	}
}
