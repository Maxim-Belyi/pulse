package redis

import (
	"context"
	"fmt"
	"pulse/internal/entity"
	"pulse/internal/usecase"
	"pulse/pkg/logger"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRepository struct {
	client *redis.Client
	logger *logger.Logger
}

func NewCacheRepository(client *redis.Client, logger *logger.Logger) *CacheRepository {
	return &CacheRepository{
		client: client,
		logger: logger,
	}
}

func (r *CacheRepository) Get(ctx context.Context, key string) ([]byte, error) {
	return r.client.Get(ctx, key).Bytes()
}

func (r *CacheRepository) Set(ctx context.Context, key string, data []byte, ttl time.Duration) error {
	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *CacheRepository) IncSourceCount(ctx context.Context, source entity.SourceType) error {
	key := fmt.Sprintf("events:count:%s:%s", source, time.Now().UTC().Format("2006-01-02"))
	return r.client.Incr(ctx, key).Err()
}

var _ usecase.CacheRepository = (*CacheRepository)(nil)