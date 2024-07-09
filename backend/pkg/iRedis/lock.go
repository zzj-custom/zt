package iRedis

import (
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
)

type Lock struct {
	Pool        *redis.Pool
	LockSeconds int
}

func NewLock(pool *redis.Pool) *Lock {
	lock := new(Lock)
	lock.Pool = pool
	lock.LockSeconds = 60
	return lock
}

func (r *Lock) Acquire(lock string, lockSeconds int) error {
	if lockSeconds <= 0 {
		lockSeconds = r.LockSeconds
	}
	conn := r.Pool.Get()
	defer func() {
		_ = conn.Close()
	}()
	_, err := redis.String(conn.Do("SET", lock, 1, "EX", lockSeconds, "NX"))
	if err != nil && err != redis.ErrNil {
		return errors.Wrapf(err, "获取锁失败，lock=%s", lock)
	}
	return nil
}

func (r *Lock) Release(lock string) error {
	conn := r.Pool.Get()
	defer func() { _ = conn.Close() }()

	_, err := conn.Do("DEL", conn)
	if err != nil {
		return errors.Wrapf(err, "删除锁失败,key=%s", lock)
	}
	return nil
}
