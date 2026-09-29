package catalog

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

const TemplateName = "atlasflow-map-catalog"

type Catalog struct {
	Version      string    `json:"version"`
	DefaultMapID any       `json:"defaultMapId"`
	Maps         []MapItem `json:"maps"`
}

type MapItem struct {
	ID          any             `json:"id"`
	Name        string          `json:"name"`
	SourceKind  string          `json:"sourceKind"`
	Source      json.RawMessage `json:"source"`
	Layers      []Layer         `json:"layers"`
	Default     bool            `json:"default"`
	Mandatory   bool            `json:"mandatory"`
	Priority    int             `json:"priority"`
	FallbackID  any             `json:"fallbackId"`
	Transport   Transport       `json:"transport"`
	Description string          `json:"description"`
}

type Layer struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Source string `json:"source"`
}

type Transport struct {
	Kind         string `json:"kind"`
	PreserveHost bool   `json:"preserveHost"`
	Encrypted    bool   `json:"encrypted"`
}

type sourceContract struct {
	Kind     string   `json:"kind"`
	Type     string   `json:"type"`
	URL      string   `json:"url"`
	Tiles    []string `json:"tiles"`
	TileSize int      `json:"tileSize"`
}

func ValidateJSON(data []byte) error {
	var value Catalog
	if err := json.Unmarshal(data, &value); err != nil {
		return fmt.Errorf("decode AtlasFlow map catalog: %w", err)
	}
	return Validate(value)
}

func Validate(value Catalog) error {
	if strings.TrimSpace(value.Version) == "" {
		return errors.New("catalog version is required")
	}
	if len(value.Maps) == 0 {
		return errors.New("catalog must contain at least one map")
	}

	seen := make(map[string]struct{}, len(value.Maps))
	fallbackRefs := make(map[string]string, len(value.Maps))
	defaultCount := 0
	defaultItemID := ""
	styleCount := 0

	for index, item := range value.Maps {
		id := strings.TrimSpace(fmt.Sprint(item.ID))
		if id == "" || id == "<nil>" {
			return fmt.Errorf("maps[%d].id is required", index)
		}
		if _, exists := seen[id]; exists {
			return fmt.Errorf("duplicate map id %q", id)
		}
		seen[id] = struct{}{}
		if strings.TrimSpace(item.Name) == "" {
			return fmt.Errorf("maps[%d].name is required", index)
		}
		if item.Default {
			defaultCount++
			defaultItemID = id
		}

		var source sourceContract
		if len(item.Source) == 0 {
			return fmt.Errorf("maps[%d].source is required", index)
		}
		if err := json.Unmarshal(item.Source, &source); err != nil {
			return fmt.Errorf("maps[%d].source is invalid: %w", index, err)
		}
		kind := strings.ToLower(strings.TrimSpace(firstNonEmpty(item.SourceKind, source.Kind)))
		switch kind {
		case "style":
			styleCount++
			if err := validateHTTPURL(source.URL); err != nil {
				return fmt.Errorf("maps[%d].source.url: %w", index, err)
			}
		case "tiles":
			if len(source.Tiles) == 0 {
				return fmt.Errorf("maps[%d].source.tiles is required", index)
			}
			for tileIndex, tileURL := range source.Tiles {
				if err := validateTileTemplate(tileURL); err != nil {
					return fmt.Errorf("maps[%d].source.tiles[%d]: %w", index, tileIndex, err)
				}
			}
		case "tilejson":
			if err := validateHTTPURL(source.URL); err != nil {
				return fmt.Errorf("maps[%d].source.url: %w", index, err)
			}
		default:
			return fmt.Errorf("maps[%d].sourceKind %q is unsupported", index, kind)
		}

		transportKind := strings.ToLower(strings.TrimSpace(item.Transport.Kind))
		switch transportKind {
		case "public-http", "ciproto-plain", "ciproto-encrypted":
		default:
			return fmt.Errorf("maps[%d].transport.kind %q is unsupported", index, item.Transport.Kind)
		}
		if transportKind == "ciproto-encrypted" && !item.Transport.Encrypted {
			return fmt.Errorf("maps[%d].transport.encrypted must be true for ciproto-encrypted", index)
		}

		fallbackID := strings.TrimSpace(fmt.Sprint(item.FallbackID))
		if fallbackID != "" && fallbackID != "<nil>" {
			if fallbackID == id {
				return fmt.Errorf("maps[%d].fallbackId cannot reference itself", index)
			}
			fallbackRefs[id] = fallbackID
		}
	}

	for id, fallbackID := range fallbackRefs {
		if _, exists := seen[fallbackID]; !exists {
			return fmt.Errorf("map %q fallbackId %q does not reference an existing map", id, fallbackID)
		}
	}

	if defaultCount != 1 {
		return fmt.Errorf("catalog must contain exactly one default map, got %d", defaultCount)
	}
	if styleCount == 0 {
		return errors.New("catalog must contain a style-owned base map")
	}
	if strings.TrimSpace(fmt.Sprint(value.DefaultMapID)) == "" {
		return errors.New("defaultMapId is required")
	}
	resolvedDefaultID := strings.TrimSpace(fmt.Sprint(value.DefaultMapID))
	if _, exists := seen[resolvedDefaultID]; !exists {
		return errors.New("defaultMapId does not reference an existing map")
	}
	if resolvedDefaultID != defaultItemID {
		return fmt.Errorf(
			"defaultMapId %q does not match the map marked default %q",
			resolvedDefaultID,
			defaultItemID,
		)
	}
	return nil
}

func validateHTTPURL(value string) error {
	parsed, err := url.Parse(strings.TrimSpace(value))
	if err != nil {
		return err
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return errors.New("must use http or https")
	}
	if parsed.Host == "" {
		return errors.New("host is required")
	}
	return nil
}

func validateTileTemplate(value string) error {
	if err := validateHTTPURL(value); err != nil {
		return err
	}
	for _, token := range []string{"{z}", "{x}", "{y}"} {
		if !strings.Contains(value, token) {
			return fmt.Errorf("missing %s placeholder", token)
		}
	}
	return nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
