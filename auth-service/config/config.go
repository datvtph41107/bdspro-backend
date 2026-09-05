package config

import (
	"common/configloader"
	"fmt"

	"github.com/spf13/viper"
)

// LoadProperties reads Auth configuration and returns startup failures to the process owner.
func LoadProperties() error {
	_, err := configloader.LoadRuntimeYML()
	if err != nil {
		return err
	}
	// YAML chỉ giữ topology/default không nhạy cảm. Secret và override của từng
	// môi trường được Viper đọc từ biến môi trường tại process boundary này.

	if err := viper.Unmarshal(&Properties); err != nil {
		return fmt.Errorf("parse auth properties: %w", err)
	}

	return nil
}
