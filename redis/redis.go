package redis

import (
	"context"
	"time"

	rd "github.com/gomodule/redigo/redis"
)

type Config struct {
	Address   string
	Password  string
	MaxIdle   int
	MaxActive int
}

type redis struct {
	pool *rd.Pool
}

func New(cfg *Config) Redis {

	pool := &rd.Pool{
		MaxIdle:     cfg.MaxIdle,
		MaxActive:   cfg.MaxActive,
		IdleTimeout: 240 * time.Second,
		Dial: func() (rd.Conn, error) {
			return rd.Dial("tcp", cfg.Address, rd.DialPassword(cfg.Password), rd.DialReadTimeout(5*time.Second), rd.DialWriteTimeout(5*time.Second))
		},
	}
	return &redis{pool}
}

type Redis interface {
	Do(ctx context.Context, commandName string, args ...interface{}) (reply interface{}, err error)
	Close() error
	Conn() rd.Conn
}

func (r *redis) Do(ctx context.Context, commandName string, args ...interface{}) (reply interface{}, err error) {
	conn, err := r.pool.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return conn.Do(commandName, args...)
}

func (r *redis) Close() error {
	return r.pool.Close()
}

func (r *redis) Conn() rd.Conn {
	return r.pool.Get()
}
