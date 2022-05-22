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
	Pipeline(ctx context.Context) Pipeline

	Del(ctx context.Context, key string) error
}

type Pipeline interface {
	Send(command string, args ...interface{})
	Receive() (replies []interface{}, err error)
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

func (r *redis) Del(ctx context.Context, key string) error {
	_, err := r.Do(ctx, "DEL", key)
	return err
}

func (r *redis) Pipeline(ctx context.Context) Pipeline {
	return &pipeline{
		conn: r.pool.Get(),
	}
}

type pipeline struct {
	conn rd.Conn
	num  int
}

func (p *pipeline) Send(command string, args ...interface{}) {
	p.conn.Send(command, args...)
	p.num++
}

func (p *pipeline) Receive() (replies []interface{}, err error) {
	defer p.conn.Close()

	if err = p.conn.Flush(); err != nil {
		return
	}

	for i := 0; i < p.num; i++ {
		reply, err := p.conn.Receive()
		if err != nil {
			return nil, err
		}
		replies = append(replies, reply)
	}
	return
}
