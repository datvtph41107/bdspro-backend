// Package logging sở hữu cơ chế xuất log; service vẫn sở hữu ý nghĩa sự kiện.
package logging

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type SQLConfig struct {
	Enabled       bool          `yaml:"enabled"`
	SlowThreshold time.Duration `yaml:"slow_threshold"`
}

type Config struct {
	Level     string `yaml:"level"`
	Console   bool   `yaml:"console"`
	Format    string `yaml:"format"`
	AddSource bool   `yaml:"add_source"`
	Files     bool   `yaml:"files"`
	Root      string `yaml:"root"`
	// session: một thư mục mỗi lần start; append: giữ file; truncate: xóa nội dung file do logger sở hữu.
	Mode       string            `yaml:"mode"`
	MaxSizeMB  int               `yaml:"max_size_mb"`
	MaxBackups int               `yaml:"max_backups"`
	MaxAgeDays int               `yaml:"max_age_days"`
	KeepRuns   int               `yaml:"keep_runs"`
	Categories map[string]string `yaml:"categories"`
	SQL        SQLConfig         `yaml:"sql"`
}

type fileConfig struct {
	Version  int                  `yaml:"version"`
	Defaults yaml.Node            `yaml:"defaults"`
	Services map[string]yaml.Node `yaml:"services"`
}

func defaults() Config {
	return Config{Level: "info", Console: true, Format: "json", Mode: "session", MaxSizeMB: 20, MaxBackups: 3, MaxAgeDays: 7, KeepRuns: 10, SQL: SQLConfig{SlowThreshold: time.Second}}
}

// Load đọc đúng một file chỉ định, không dò tìm config theo cwd hay âm thầm dùng file khác.
// Không chỉ định file: stdout JSON phù hợp container; service override được merge lên defaults.
func Load(path, service string) (Config, error) {
	c := defaults()
	if path == "" {
		return c, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return c, err
	}
	defer f.Close()
	var document fileConfig
	decoder := yaml.NewDecoder(f)
	decoder.KnownFields(true)
	if err = decoder.Decode(&document); err != nil {
		return c, err
	}
	if document.Version != 1 {
		return c, errors.New("logging config version must be 1")
	}
	// Decode từng lớp với KnownFields để typo không làm dev tưởng cấu hình đã có hiệu lực.
	for _, node := range []yaml.Node{document.Defaults, document.Services[service]} {
		if node.Kind == 0 {
			continue
		}
		data, err := yaml.Marshal(node)
		if err != nil {
			return c, err
		}
		d := yaml.NewDecoder(strings.NewReader(string(data)))
		d.KnownFields(true)
		if err = d.Decode(&c); err != nil {
			return c, err
		}
	}
	if c.Root != "" && !filepath.IsAbs(c.Root) {
		c.Root = filepath.Join(filepath.Dir(path), c.Root)
	}
	return c, c.validate()
}

func (c Config) validate() error {
	if _, err := parseLevel(c.Level); err != nil {
		return err
	}
	if c.Format != "json" && c.Format != "text" {
		return errors.New("logging format must be json or text")
	}
	if c.Mode != "session" && c.Mode != "append" && c.Mode != "truncate" {
		return errors.New("logging mode must be session, append or truncate")
	}
	if c.Files && (c.Root == "" || c.MaxSizeMB <= 0 || c.MaxBackups <= 0 || c.MaxAgeDays <= 0 || c.KeepRuns <= 0) {
		return errors.New("file logging requires root and positive rotation limits")
	}
	if c.SQL.SlowThreshold <= 0 {
		return errors.New("sql slow_threshold must be positive")
	}
	for category, level := range c.Categories {
		if !validCategory(category) {
			return fmt.Errorf("unknown log category %q", category)
		}
		if level != "off" {
			if _, err := parseLevel(level); err != nil {
				return err
			}
		}
	}
	return nil
}

var categories = []string{"application", "http", "grpc", "database", "dependency", "worker", "lifecycle"}

func validCategory(v string) bool {
	for _, c := range categories {
		if c == v {
			return true
		}
	}
	return false
}
func parseLevel(v string) (slog.Level, error) {
	var l slog.Level
	err := l.UnmarshalText([]byte(v))
	return l, err
}
