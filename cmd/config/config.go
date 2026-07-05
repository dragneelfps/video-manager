package config

import (
	"fmt"

	"github.com/spf13/viper"
	"nooblabs.com/video-manager/internal/model"
)

func LoadConfig(v *viper.Viper) (model.Config, error) {
	var cfg model.Config
	err := v.Unmarshal(&cfg)
	if err != nil {
		return model.Config{}, fmt.Errorf("failed to read config - %w", err)
	}
	return cfg, nil
}
