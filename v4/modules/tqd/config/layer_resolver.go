package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"mdn4/modules/tqd/enums"
)

// LayerConfig holds tile/style settings used by LayerResolver (ported from tqd-service/config).
type LayerConfig struct {
	BaseURL       string                   `mapstructure:"base_url"`
	Overrides     map[string]LayerOverride `mapstructure:"overrides"`
	DefaultStyles LayerDefaultStyles       `mapstructure:"default_styles"`
}

type LayerOverride struct {
	URLPattern  string `mapstructure:"url_pattern"`
	StyleConfig string `mapstructure:"style_config"`
}

type LayerDefaultStyles struct {
	Vector []map[string]interface{} `mapstructure:"vector"`
}

type LayerResolver struct {
	config *LayerConfig
}

// LayerResolverInstance is set by InitLayerResolver (or left as a safe empty resolver).
var LayerResolverInstance = NewLayerResolver(&LayerConfig{
	BaseURL:   "",
	Overrides: map[string]LayerOverride{},
	DefaultStyles: LayerDefaultStyles{
		Vector: []map[string]interface{}{},
	},
})

func NewLayerResolver(cfg *LayerConfig) *LayerResolver {
	if cfg == nil {
		cfg = &LayerConfig{Overrides: map[string]LayerOverride{}}
	}
	if cfg.Overrides == nil {
		cfg.Overrides = map[string]LayerOverride{}
	}
	return &LayerResolver{config: cfg}
}

// InitLayerResolver replaces the process-wide resolver (call after loading runtime config).
func InitLayerResolver(cfg *LayerConfig) {
	LayerResolverInstance = NewLayerResolver(cfg)
}

func (r *LayerResolver) GetLayerURL(layerID uint64, sourceCode string) string {
	layerIDStr := fmt.Sprintf("%d", layerID)
	if override, exists := r.config.Overrides[layerIDStr]; exists && override.URLPattern != "" {
		return r.config.BaseURL + override.URLPattern
	}
	return fmt.Sprintf("%s/%s/{z}/{x}/{y}.pbf", r.config.BaseURL, sourceCode)
}

func (r *LayerResolver) GetStyleConfig(layerID uint64, sourceCode string, layerType uint32, minZoom uint32, maxZoom uint32) (string, error) {
	layerIDStr := fmt.Sprintf("%d", layerID)
	if override, exists := r.config.Overrides[layerIDStr]; exists && override.StyleConfig != "" {
		return override.StyleConfig, nil
	}
	_ = enums.LayerType(layerType)
	defaultStyle := r.config.DefaultStyles.Vector
	styleJSON, err := json.Marshal(defaultStyle)
	if err != nil {
		return "[]", err
	}
	styleStr := string(styleJSON)
	styleStr = strings.ReplaceAll(styleStr, "{sourceCode}", sourceCode)
	styleStr = strings.ReplaceAll(styleStr, "{maxZoom}", fmt.Sprintf("%d", maxZoom))
	styleStr = strings.ReplaceAll(styleStr, "{minZoom}", fmt.Sprintf("%d", minZoom))
	return styleStr, nil
}

func (r *LayerResolver) GetSourceCode(layerID uint64) string {
	return fmt.Sprintf("layer%d", layerID)
}

func (r *LayerResolver) GetTileName(layerID uint64) string {
	return fmt.Sprintf("ci_layer%d", layerID)
}

func (r *LayerResolver) GetFamilyTileName(familyID uint64) string {
	return fmt.Sprintf("ci_family%d", familyID)
}

func (r *LayerResolver) GetFamilyURL(familyID uint64, tileName string) string {
	return fmt.Sprintf("%s/%s/{z}/{x}/{y}.pbf", r.config.BaseURL, tileName)
}

func (r *LayerResolver) GetFamilySourceCode(familyID uint64) string {
	return fmt.Sprintf("family%d", familyID)
}

func (r *LayerResolver) GetFamilyStyleConfig(familyID uint64, minZoom, maxZoom uint32) (string, error) {
	familyIDStr := fmt.Sprintf("family_%d", familyID)
	if override, exists := r.config.Overrides[familyIDStr]; exists && override.StyleConfig != "" {
		return override.StyleConfig, nil
	}
	return r.GetStyleConfig(0, r.GetFamilySourceCode(familyID), uint32(enums.LayerTypeLandUse), minZoom, maxZoom)
}
