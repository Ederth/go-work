package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Hash struct {
	Client   *redis.Client
	RedisKey string
}

func (h *Hash) HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	return h.Client.HSet(ctx, h.getKey(key), values...)
}

func (h *Hash) HGet(ctx context.Context, key, field string) *redis.StringCmd {
	return h.Client.HGet(ctx, h.getKey(key), field)
}

func (h *Hash) HMGet(ctx context.Context, key string, fields ...string) *redis.SliceCmd {
	return h.Client.HMGet(ctx, h.getKey(key), fields...)
}

func (h *Hash) HGetAll(ctx context.Context, key string) *redis.StringStringMapCmd {
	return h.Client.HGetAll(ctx, h.getKey(key))
}

func (h *Hash) getKey(key string) string {
	if key != "" {
		return h.RedisKey + ":" + key
	}

	return h.RedisKey
}
