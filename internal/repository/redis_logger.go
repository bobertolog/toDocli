package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisLogger struct {
	client *redis.Client
}

func NewRedisLogger(c *redis.Client) *RedisLogger {
	return &RedisLogger{client: c}
}

func (l *RedisLogger) Log(msg string) error {
	key := "log:" + uuid.NewString()
	return l.client.Set(context.Background(), key, msg, time.Hour).Err()
}
