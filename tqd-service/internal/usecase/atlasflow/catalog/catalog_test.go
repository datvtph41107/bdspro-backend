package catalog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAtlasFlowCatalogTemplateIsValid(t *testing.T) {
	path := filepath.Join("..", "..", "..", "..", "template", TemplateName+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read catalog template: %v", err)
	}

	content := strings.ReplaceAll(string(data), "{{domain_api}}", "https://api.example.test")
	if err := ValidateJSON([]byte(content)); err != nil {
		t.Fatalf("catalog template invalid: %v", err)
	}
}

func TestCatalogRejectsDuplicateIDs(t *testing.T) {
	data := []byte(`{
		"version":"2.0.0",
		"defaultMapId":"street",
		"maps":[
			{"id":"street","name":"One","sourceKind":"style","source":{"kind":"style","url":"https://api.example/style"},"default":true,"transport":{"kind":"public-http"}},
			{"id":"street","name":"Two","sourceKind":"style","source":{"kind":"style","url":"https://api.example/style"},"transport":{"kind":"public-http"}}
		]
	}`)
	if err := ValidateJSON(data); err == nil {
		t.Fatal("expected duplicate id validation error")
	}
}

func TestCatalogRejectsUnknownFallback(t *testing.T) {
	data := []byte(`{
		"version":"2.0.0",
		"defaultMapId":"street",
		"maps":[
			{"id":"street","name":"One","sourceKind":"style","source":{"kind":"style","url":"https://api.example/style"},"default":true,"transport":{"kind":"public-http"}},
			{"id":"satellite","name":"Two","sourceKind":"tiles","source":{"kind":"tiles","type":"raster","tiles":["https://tiles.example/{z}/{x}/{y}.png"]},"fallbackId":"missing","transport":{"kind":"public-http"}}
		]
	}`)
	if err := ValidateJSON(data); err == nil {
		t.Fatal("expected unknown fallback validation error")
	}
}

func TestCatalogRejectsInvalidEncryptedTransport(t *testing.T) {
	data := []byte(`{
		"version":"2.0.0",
		"defaultMapId":"street",
		"maps":[
			{"id":"street","name":"One","sourceKind":"style","source":{"kind":"style","url":"https://api.example/style"},"default":true,"transport":{"kind":"ciproto-encrypted","encrypted":false}}
		]
	}`)
	if err := ValidateJSON(data); err == nil {
		t.Fatal("expected encrypted transport validation error")
	}
}

func TestCatalogRejectsDefaultMapMismatch(t *testing.T) {
	data := []byte(`{
		"version":"2.0.0",
		"defaultMapId":"satellite",
		"maps":[
			{"id":"street","name":"One","sourceKind":"style","source":{"kind":"style","url":"https://api.example/style"},"default":true,"transport":{"kind":"public-http"}},
			{"id":"satellite","name":"Two","sourceKind":"tiles","source":{"kind":"tiles","type":"raster","tiles":["https://tiles.example/{z}/{x}/{y}.png"]},"transport":{"kind":"public-http"}}
		]
	}`)
	if err := ValidateJSON(data); err == nil {
		t.Fatal("expected defaultMapId mismatch validation error")
	}
}
