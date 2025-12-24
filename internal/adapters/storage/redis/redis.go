package redis

import (
	"context"
	"fmt"

	"app/config"
	"app/pkg/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(cfg config.RedisConfig) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	logger.Info("Connected to Redis", zap.String("addr", cfg.Addr))
	return &RedisClient{Client: rdb}, nil
}

func (r *RedisClient) Close() {
	if r.Client != nil {
		r.Client.Close()
	}
}
