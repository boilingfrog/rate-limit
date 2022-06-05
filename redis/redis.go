package redis

import (
	rd "github.com/go-redis/redis/v8"
)

type Config struct {
	Address  string
	Password string
}

type Redis struct {
	*rd.Client
}

func New(cfg *Config) *Redis {

	rdb := rd.NewClient(&rd.Options{
		Addr:     cfg.Address,
		Password: cfg.Password, // no password set
		DB:       0,            // use default DB
	})

	return &Redis{rdb}
}
