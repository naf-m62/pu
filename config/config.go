package config

import (
	"github.com/spf13/viper"
)

type Config = *viper.Viper

func New(configPath string) (config Config, err error) {
	var cfg = viper.New()
	cfg.SetConfigFile(configPath)

	cfg.SetConfigType("yaml")

	if err = cfg.ReadInConfig(); err != nil {
		return nil, err
	}

	return cfg, nil
}
