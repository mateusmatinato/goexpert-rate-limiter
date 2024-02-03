package redis

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addr     string
	Port     int
	Password string
}

func NewClient(cfg Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port),
		Password: cfg.Password,
		DB:       0,
	})
	return client
}
