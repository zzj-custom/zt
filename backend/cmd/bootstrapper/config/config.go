package config

import (
	"github.com/spf13/viper"
)

var global = "global"

func Store(value any) {
	viper.Set(global, value)
}

func Get() *Config {
	cfg := viper.Get(global)
	return cfg.(*Config)
}
