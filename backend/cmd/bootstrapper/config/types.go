package config

import (
	"zt/backend/pkg/email"
	"zt/backend/pkg/iMysql"
	"zt/backend/pkg/iRedis"
	"zt/backend/pkg/sms"
)

type Config struct {
	Database map[string]*iMysql.Database `toml:"database"`
	Redis    []*iRedis.MultiDialConfig   `toml:"redis"`
	Email    *email.Config               `toml:"email"`
	Sms      *sms.Config                 `toml:"sms"`
}
