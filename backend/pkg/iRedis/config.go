package iRedis

import (
	"github.com/wbylovesun/xutils/xserializer"
)

type DialConfig struct {
	Host            string               `toml:"host"`
	Port            int                  `toml:"port"`
	Database        int                  `toml:"database"`
	Password        string               `toml:"password"`
	MaxIdle         int                  `toml:"max_idle"`
	MaxActive       int                  `toml:"max_active"`
	Wait            bool                 `toml:"wait"`
	ConnectTimeout  xserializer.Duration `toml:"connect_timeout"`
	ReadTimeout     xserializer.Duration `toml:"read_timeout"`
	MaxConnLifetime xserializer.Duration `toml:"max_conn_lifetime"`
	IdleTimeout     xserializer.Duration `toml:"idle_timeout"`
}

// MultiDialConfig
//
// Define MultiDialConfig item with name surrounded by `[[` and `]]`.
//
// Example:
// [[redis]]
// name="first"
// default=true
// config.host="localhost"
// config.port=6379
// config.database=0
// config.connect_timeout="5s"
// config.read_timeout="2s"
// config.max_idle=1
// config.max_active=3
// config.idle_timeout="60s"
// config.wait=true
// config.max_conn_lifetime="3600s"
// [[redis]]
// name="second"
// config.host="localhost"
// config.port=6379
// config.database=1
// config.connect_timeout="5s"
// config.read_timeout="2s"
// config.max_idle=1
// config.max_active=3
// config.idle_timeout="60s"
// config.wait=true
// config.max_conn_lifetime="3600s"
type MultiDialConfig struct {
	Name    string      `toml:"name"`
	Default bool        `toml:"default"`
	Config  *DialConfig `toml:"config"`
}
