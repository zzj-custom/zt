package bootstrapper

import (
	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"log/slog"
	"zt/backend/cmd/bootstrapper/config"
)

func init() {
	viper.SetDefault("config.name", "config")
	viper.SetDefault("config.type", "toml")
}

func Bootstrap(cfgFile string, f func(in fsnotify.Event)) *config.Config {
	// 加载配置文件
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(viper.GetString("config.name"))
		viper.SetConfigType(viper.GetString("config.type"))
	}
	if err := viper.ReadInConfig(); err != nil {
		slog.With("err", err).Error("read config file error")
		panic(err)
	}

	// 加载配置
	var v config.Config
	bootstrap(&v, f)
	config.Store(&v)
	// 初始化资源
	initResource()

	return &v
}

func bootstrap(v *config.Config, f func(in fsnotify.Event)) {
	viper.OnConfigChange(f)
	viper.WatchConfig()
	err := viper.Unmarshal(&v, viper.DecodeHook(mapstructure.StringToTimeDurationHookFunc()))
	if err != nil {
		slog.With("err", err).Error("unmarshal config error")
		panic(err)
	}
}
