package config

import (
	"common/configloader"
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
	"tqd/internal/enums"

	"github.com/spf13/viper"
)

// Properties holds all configuration properties
type Properties struct {
	Server    ServerConfig   `mapstructure:"server"`
	Database  DatabaseConfig `mapstructure:"database"`
	JWT       JWTConfig      `mapstructure:"jwt"`
	Eureka    EurekaConfig   `mapstructure:"eureka"`
	Key       KeyConfig      `mapstructure:"key"`
	Layer     LayerConfig    `mapstructure:"layer"`
	ConfigApp struct {
		Qhpro struct {
			DomainAPI    string `mapstructure:"domain_api"`
			DomainFile   string `mapstructure:"domain_file"`
			DomainPMTile string `mapstructure:"domain_pmtile"`
			DomainStyle  string `mapstructure:"domain_style"`
		} `mapstructure:"qhpro"`
	} `mapstructure:"config-app"`
}

type LayerConfig struct {
	BaseURL       string                   `mapstructure:"base_url"`
	DefaultStyles DefaultStylesConfig      `mapstructure:"default_styles"`
	Overrides     map[string]LayerOverride `mapstructure:"overrides"`
}

type DefaultStylesConfig struct {
	Vector []map[string]interface{} `mapstructure:"vector"`
	Raster []map[string]interface{} `mapstructure:"raster"`
}

type LayerOverride struct {
	URLPattern  string `mapstructure:"url_pattern"`
	StyleConfig string `mapstructure:"style_config"`
}

type LayerResolver struct {
	config *LayerConfig
}

func NewLayerResolver() *LayerResolver {
	return &LayerResolver{
		config: &AppConfig.Layer,
	}
}

// GetLayerURL - Lấy URL cho layer theo id và sourceCode
func (r *LayerResolver) GetLayerURL(layerID uint64, sourceCode string) string {
	layerIDStr := fmt.Sprintf("%d", layerID)
	// Kiểm tra override
	if override, exists := r.config.Overrides[layerIDStr]; exists && override.URLPattern != "" {
		return r.config.BaseURL + override.URLPattern
	}

	// Default pattern: /{sourceCode}/{z}/{x}/{y}.pbf
	return fmt.Sprintf("%s/%s/{z}/{x}/{y}.pbf", r.config.BaseURL, sourceCode)
}

// GetStyleConfig - Lấy style config cho layer
func (r *LayerResolver) GetStyleConfig(layerID uint64, sourceCode string, layerType uint32, minZoom uint32, maxZoom uint32) (string, error) {
	layerIDStr := fmt.Sprintf("%d", layerID)

	// override config
	if override, exists := r.config.Overrides[layerIDStr]; exists && override.StyleConfig != "" {
		return override.StyleConfig, nil
	}

	// Lấy default style theo type
	var defaultStyle []map[string]interface{}

	// LayerType: 1=landuse, 3=planning, 4=traffic
	switch enums.LayerType(layerType) {
	// case enums.LayerTypeLandUse, enums.LayerTypePlanning:
	// 	defaultStyle = r.config.DefaultStyles.Vector
	// case enums.LayerTypeTraffic:
	// 	defaultStyle = r.config.DefaultStyles.Vector
	default:
		defaultStyle = r.config.DefaultStyles.Vector
	}

	// Replace placeholders
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

// GetFamilyURL - Lấy URL tile cho họ lớp theo id
func (r *LayerResolver) GetFamilyURL(familyID uint64, tileName string) string {
	return fmt.Sprintf("%s/%s/{z}/{x}/{y}.pbf", r.config.BaseURL, tileName)
}

func (r *LayerResolver) GetFamilySourceCode(familyID uint64) string {
	return fmt.Sprintf("family%d", familyID)
}

// GetFamilyStyleConfig - Lấy style config cho họ lớp (cùng template vector với layer)
func (r *LayerResolver) GetFamilyStyleConfig(familyID uint64, minZoom, maxZoom uint32) (string, error) {
	familyIDStr := fmt.Sprintf("family_%d", familyID)
	if override, exists := r.config.Overrides[familyIDStr]; exists && override.StyleConfig != "" {
		return override.StyleConfig, nil
	}
	return r.GetStyleConfig(0, r.GetFamilySourceCode(familyID), uint32(enums.LayerTypeLandUse), minZoom, maxZoom)
}

// ServerConfig holds server configuration
type ServerConfig struct {
	GrpcPort int `mapstructure:"tcp_port"`
	HttpPort int `mapstructure:"port"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	DSN string `mapstructure:"dsn"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	KeyGenerate         string   `mapstructure:"key-generate"`
	TokenPrefix         string   `mapstructure:"tokenPrefix"`
	TokenExpirationDays int      `mapstructure:"tokenExpirationAfterDays"`
	AccessExpMinutes    int      `mapstructure:"accessExpAfterMinutes"`
	RefreshExpMinutes   int      `mapstructure:"refreshExpAfterMinutes"`
	AuthorizationHeader string   `mapstructure:"authorizationHeader"`
	ListPermit          []string `mapstructure:"listPermit"`
}

// EurekaConfig holds Eureka configuration
type EurekaConfig struct {
	URL  string `mapstructure:"url"`
	Port int    `mapstructure:"port"`
}

// KeyConfig holds key configuration
type KeyConfig struct {
	JWT string `mapstructure:"jwt"`
}

var AppConfig *Properties
var LayerResolverInstance *LayerResolver

// LoadConfig đọc cấu hình từ file config.yml
func LoadConfig() (*Properties, error) {
	_, err := configloader.LoadRuntimeYML()
	if err != nil {
		return nil, err
	}

	AppConfig = &Properties{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		return nil, fmt.Errorf("unmarshal tqd config: %w", err)
	}

	LayerResolverInstance = NewLayerResolver()
	return AppConfig, nil
}
