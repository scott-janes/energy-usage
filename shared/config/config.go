package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func LoadConfig(configPath string, configName string, configType string, config interface{}) error {
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(configPath)

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("Error reading config file: %w", err)
	}

	err = viper.Unmarshal(config)

	if err != nil {
		return fmt.Errorf("Error unmarshaling config: %w", err)
  }

	return nil
}
