package access

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go -package=mocks
type Repository interface {
	GetAccessCount(ctx context.Context, key string) (int, error)
	IncrementAccessCount(ctx context.Context, key string) error
}

type repository struct {
	redisClient *redis.Client
}

func (r repository) GetAccessCount(ctx context.Context, key string) (int, error) {
	countStr, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, err
	}
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r repository) IncrementAccessCount(ctx context.Context, key string) error {
	count, err := r.GetAccessCount(ctx, key)
	if err != nil {
		return err
	}

	expiration := 1 * time.Second
	if count > 0 {
		expiration = redis.KeepTTL
	}

	return r.redisClient.Set(ctx, key, count+1, expiration).Err()
}

func NewRepository(redisClient *redis.Client) Repository {
	return &repository{
		redisClient: redisClient,
	}
}
