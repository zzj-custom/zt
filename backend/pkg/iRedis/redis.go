package iRedis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"strconv"
	"sync"
	"time"
)

var (
	redisPool       = sync.Map{}
	defaultPoolName = "default"
)

func (config *DialConfig) getDialOption() []redis.DialOption {
	dialOptions := []redis.DialOption{
		redis.DialReadTimeout(config.ReadTimeout.Duration),
		redis.DialConnectTimeout(config.ConnectTimeout.Duration),
		redis.DialDatabase(config.Database),
	}
	if config.Password != "" {
		dialOptions = append(dialOptions, redis.DialPassword(config.Password))
	}
	return dialOptions
}

// NewPool
//
// Just create a redis connection pool, you should persist it by yourself.
// RegisterPool can be used to persist it.
func NewPool(config *DialConfig) (*redis.Pool, error) {
	if config == nil {
		return nil, fmt.Errorf("invalid initializer provided")
	}

	pool := redis.Pool{
		Dial: func() (redis.Conn, error) {
			dial, err := redis.Dial(
				"tcp",
				config.Host+":"+strconv.Itoa(config.Port),
				config.getDialOption()...,
			)
			if err != nil {
				return nil, err
			}
			_, err = dial.Do("SELECT", config.Database)
			if err != nil {
				dial.Close()
				return nil, err
			}
			return dial, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:         config.MaxIdle,
		MaxActive:       config.MaxActive,
		IdleTimeout:     config.IdleTimeout.Duration,
		Wait:            config.Wait,
		MaxConnLifetime: config.MaxConnLifetime.Duration,
	}
	return &pool, nil
}

// RegisterPool
//
// Persist a redis connection pool with given name.
// If this name exists, pool will be replaced by this new one, and
// the old pool will be closed.
func RegisterPool(name string, pool *redis.Pool) {
	if p, ok := redisPool.Load(name); ok {
		p2 := p.(*redis.Pool)
		_ = p2.Close()
	}
	redisPool.Store(name, pool)
}

// UseAsDefaultPool
//
// Change the `defaultPoolName` value. Use this func with InitPool to call
// Pool directly without parameter name if parameter name of InitPool is
// not "default".
func UseAsDefaultPool(name string) {
	defaultPoolName = name
}

// InitPool
// Init redis connection pool with given name, so can fetch it again
// by this give name through `iRedis.Pool(name string)`
func InitPool(name string, config *DialConfig) (*redis.Pool, error) {
	pool, err := NewPool(config)
	if err != nil {
		return nil, err
	}
	RegisterPool(name, pool)
	return pool, nil
}

// InitDefaultPool
//
// Init redis connection pool with `defaultPoolName` value, so can fetch it again
// directly through `iRedis.Pool()` without parameter.
//
// Only can access it by `defaultPoolName` value, after this func called, UseAsDefaultPool
// no longer takes effect. And NOT SUGGEST to use it with `UseAsDefaultPool`.
//
// If you want to build multiple connection pools, you should use InitMultiPools instead.
func InitDefaultPool(config *DialConfig) (*redis.Pool, error) {
	return InitPool(defaultPoolName, config)
}

// InitMultiPools
//
// Init multiple redis connection pools. You can specify name for each redis connection pool.
// And if you want to use one pool directly, you should set `Default` to be true.
//
// If `Default` that value equals true are found more than 1 time, only the last one is taken.
// And its name will be set as `defaultPoolName`.
func InitMultiPools(configs []*MultiDialConfig) error {
	if len(configs) == 0 {
		return errors.Errorf("no valid configs specified")
	}
	hasDefaultFound := false
	for _, mdc := range configs {
		_, err := InitPool(mdc.Name, mdc.Config)
		if err != nil {
			return err
		}
		if mdc.Default {
			UseAsDefaultPool(mdc.Name)
			hasDefaultFound = true
		}
	}
	if !hasDefaultFound {
		defaultPoolName = configs[0].Name
	}
	return nil
}

// Pool
//
// Fetch a redis connection pool with given name. If no name given, use
// `defaultPoolName` to fetch the default pool.
func Pool(name ...string) (*redis.Pool, error) {
	poolName := defaultPoolName
	if len(name) > 0 {
		poolName = name[0]
	}
	p, ok := redisPool.Load(poolName)
	if !ok {
		return nil, errors.Errorf("no such pool %s found", poolName)
	}
	pool := p.(*redis.Pool)
	return pool, nil
}

// Release 释放所有连接池
func Release() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("iRedis Release panic:", err)
		}
	}()
	redisPool.Range(func(key, value any) bool {
		pool := value.(*redis.Pool)
		if pool != nil && pool.Get() != nil {
			_ = pool.Close()
		}
		return true
	})
}
