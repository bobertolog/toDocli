package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLogger struct {
	client *redis.Client
}

func NewRedisLogger(client *redis.Client) *RedisLogger {
	return &RedisLogger{client: client}
}

func (r *RedisLogger) Log(event string, data string) error {
	ctx := context.TODO()
	key := fmt.Sprintf("log:%s:%d", event, time.Now().UnixNano())
	return r.client.Set(ctx, key, data, 5*time.Minute).Err()
}
