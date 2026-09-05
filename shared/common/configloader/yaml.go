// Package configloader owns process-neutral configuration loading mechanisms.
// It returns errors to the composition owner and never decides process exit.
package configloader

import (
	"bytes"
	_config "common/config"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func configFilePath(runtimeEnv string) string {
	env := strings.TrimSpace(runtimeEnv)
	if env == "" {
		env = "local"
	}
	return "config/" + env + ".yml"
}

func configureViper() {
	viper.SetConfigType("yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	// Existing SEO defaults are part of the current runtime contract. This
	// migration changes error ownership only; it does not redefine business
	// configuration defaults.
	viper.SetDefault("seo.public_base_url", "http://localhost:3000")
	viper.SetDefault("seo.internal_secret_key", "")
	viper.SetDefault("seo.timezone", "Asia/Ho_Chi_Minh")
	viper.SetDefault("seo.sitemap.default_change_freq", "daily")
	viper.SetDefault("seo.sitemap.default_priority", 0.5)
	viper.SetDefault("seo.generation_job.enabled", true)
	viper.SetDefault("seo.generation_job.cron", "0 0 * * *")
	viper.SetDefault("seo.generation_job.batch_size", 200)
}

// LoadYMLFile loads config/<runtime>.yml into the shared Viper registry.
// Empty runtime names preserve the historical local.yml behavior.
func LoadYMLFile(runtimeEnv string) error {
	configureViper()

	// Operator/deployment có thể mount một file ngoài repository. Đây là
	// override duy nhất cho production secrets/topology; nếu không có thì
	// profile canonical config/<environment>.yml được dùng.
	path := strings.TrimSpace(os.Getenv("QHPRO_CONFIG_FILE"))
	if path == "" {
		path = strings.TrimSpace(os.Getenv("CONFIG_FILE"))
	}
	if path == "" {
		path = configFilePath(runtimeEnv)
	}
	configFile, err := _config.ReadFile(path)
	if err != nil {
		return err
	}
	if err := viper.ReadConfig(bytes.NewReader(configFile)); err != nil {
		return fmt.Errorf("parse config %s: %w", path, err)
	}
	return nil
}
