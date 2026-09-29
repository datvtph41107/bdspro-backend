// Package config owns small, process-neutral configuration decoding mechanics.
// It never chooses an environment, reads process environment variables, installs
// business defaults, mutates a global registry, or decides process exit.
package config

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// ReadFile returns configuration bytes from an explicit path.
func ReadFile(path string) ([]byte, error) {
	if path == "" {
		return nil, fmt.Errorf("config path is empty")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}
	return data, nil
}

// DecodeYAML strictly decodes exactly one YAML document into a caller-owned typed value.
// Unknown fields are rejected so configuration cannot silently advertise behavior the
// running process does not own.
func DecodeYAML(data []byte, out any) error {
	if out == nil {
		return fmt.Errorf("config output is nil")
	}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(out); err != nil {
		return fmt.Errorf("parse yaml config: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return fmt.Errorf("parse yaml config: multiple YAML documents are not allowed")
		}
		return fmt.Errorf("parse yaml config: %w", err)
	}
	return nil
}

// LoadYAML reads and decodes an explicit YAML file into a caller-owned value.
func LoadYAML(path string, out any) error {
	data, err := ReadFile(path)
	if err != nil {
		return err
	}
	if err := DecodeYAML(data, out); err != nil {
		return fmt.Errorf("parse config %s: %w", path, err)
	}
	return nil
}
