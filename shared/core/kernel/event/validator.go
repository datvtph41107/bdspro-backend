package event

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/xeipuuv/gojsonschema"
)

type SchemaRegistry struct {
	dir   string
	cache map[string]*gojsonschema.Schema
	mu    sync.RWMutex
}

func NewSchemaRegistry(dir string) *SchemaRegistry {
	return &SchemaRegistry{dir: dir, cache: map[string]*gojsonschema.Schema{}}
}

func (s *SchemaRegistry) Load(eventType string, version int) (*gojsonschema.Schema, error) {
	key := fmt.Sprintf("%s.v%d", eventType, version)
	s.mu.RLock()
	if sc, ok := s.cache[key]; ok {
		s.mu.RUnlock()
		return sc, nil
	}
	s.mu.RUnlock()

	path := filepath.Join(s.dir, fmt.Sprintf("%s.v%d.json", eventType, version))
	loader := gojsonschema.NewReferenceLoader("file://" + path)
	schema, err := gojsonschema.NewSchema(loader)
	if err != nil {
		return nil, err
	}
	s.mu.Lock()
	s.cache[key] = schema
	s.mu.Unlock()
	return schema, nil
}

func (s *SchemaRegistry) Validate(eventType string, version int, payload []byte) error {
	schema, err := s.Load(eventType, version)
	if err != nil {
		return err
	}
	doc := gojsonschema.NewBytesLoader(payload)
	res, err := schema.Validate(doc)
	if err != nil {
		return err
	}
	if !res.Valid() {
		return fmt.Errorf("schema validation failed: %v", res.Errors())
	}
	return nil
}
