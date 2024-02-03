package blocked

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go -package=mocks
type Repository interface {
	Block(ctx context.Context, key string, expiration time.Duration) error
	IsBlocked(ctx context.Context, key string) (bool, error)
}

type repository struct {
	redisClient *redis.Client
}

func (r repository) Block(ctx context.Context, key string, expiration time.Duration) error {
	return r.redisClient.Set(ctx, key, "blocked", expiration).Err()
}

func (r repository) IsBlocked(ctx context.Context, key string) (bool, error) {
	res, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}

	if res != "blocked" {
		return false, nil
	}

	return true, nil
}

func NewRepository(redisClient *redis.Client) Repository {
	return &repository{
		redisClient: redisClient,
	}
}
