package code

import (
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"zt/backend/pkg/iRedis"
)

const (
	defaultExpired = 120
)

type CaptchaCache struct {
	email          string
	defaultExpired int
}

type CaptchaOption func(*CaptchaCache)

func WithCaptchaOptions(expired int) CaptchaOption {
	return func(options *CaptchaCache) {
		if expired > 0 {
			options.defaultExpired = expired
		}
	}
}

func NewCode(email string, opts ...CaptchaOption) *CaptchaCache {
	options := &CaptchaCache{
		email:          email,
		defaultExpired: defaultExpired,
	}
	for _, opt := range opts {
		opt(options)
	}

	return options
}

func (receiver CaptchaCache) key() string {
	return "captcha:" + receiver.email
}

func (receiver CaptchaCache) Validate() (bool, error) {
	pool, err := iRedis.Pool()
	if err != nil {
		return false, errors.Wrapf(err, "获取redis连接池失败")
	}
	conn := pool.Get()
	defer func() { _ = conn.Close() }()

	ok, err := redis.Bool(conn.Do("EXISTS", receiver.key()))
	if err != nil {
		return false, errors.Wrapf(err, "判断验证码是否存在失败")
	}
	return ok, nil
}

func (receiver CaptchaCache) Set(captcha int) error {
	pool, err := iRedis.Pool()
	if err != nil {
		return errors.Wrapf(err, "获取redis连接池失败")
	}
	conn := pool.Get()
	defer func() { _ = conn.Close() }()

	_, err = conn.Do("SETEX", receiver.key(), 120, captcha)
	if err != nil {
		return errors.Wrapf(err, "设置验证码失败")
	}
	return nil
}
