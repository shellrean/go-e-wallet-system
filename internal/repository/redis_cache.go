package repository

import (
	"context"
	"github.com/redis/go-redis/v9"
	"shellrean.id/belajar-auth/domain"
	"shellrean.id/belajar-auth/internal/config"
	"time"
)

type redisCacheRepository struct {
	rdb *redis.Client
}

func NewRedisClient(cnf *config.Config) domain.CacheRepository {
	return &redisCacheRepository{
		rdb: redis.NewClient(&redis.Options{
			Addr:     cnf.Redis.Addr,
			Password: cnf.Redis.Pass,
			DB:       0,
		}),
	}
}

func (r redisCacheRepository) Get(key string) ([]byte, error) {
	val, err := r.rdb.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}
	return []byte(val), nil
}

func (r redisCacheRepository) Set(key string, entry []byte) error {
	return r.rdb.Set(context.Background(), key, entry, 15*time.Minute).Err()
}
