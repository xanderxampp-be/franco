package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"

	apmgoredis "go.elastic.co/apm/module/apmgoredisv8"
)

type RedisBrimo struct {
	redis *redis.Client
}

func New(addr string, pass string) *RedisBrimo {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
	})

	client.AddHook(apmgoredis.NewHook())

	redisBrimo := &RedisBrimo{
		redis: client,
	}

	return redisBrimo
}

func (r *RedisBrimo) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return r.redis.Set(ctx, key, value, expiration)
}

func (r *RedisBrimo) Get(ctx context.Context, key string) *redis.StringCmd {
	return r.redis.Get(ctx, key)
}

func (r *RedisBrimo) HMSet(ctx context.Context, key string, values ...interface{}) *redis.BoolCmd {
	return r.redis.HMSet(ctx, key, values...)
}

func (r *RedisBrimo) HMGet(ctx context.Context, key string, fields ...string) *redis.SliceCmd {
	return r.redis.HMGet(ctx, key, fields...)
}

func (r *RedisBrimo) HIncrBy(ctx context.Context, key string, field string, incr int64) *redis.IntCmd {
	return r.redis.HIncrBy(ctx, key, field, incr)
}

func (r *RedisBrimo) HGetAll(ctx context.Context, key string) *redis.StringStringMapCmd {
	return r.redis.HGetAll(ctx, key)
}

func (r *RedisBrimo) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return r.redis.Del(ctx, keys...)
}

func (r *RedisBrimo) TTL(ctx context.Context, key string) *redis.DurationCmd {
	return r.redis.TTL(ctx, key)
}

func (r *RedisBrimo) Incr(ctx context.Context, key string) *redis.IntCmd {
	return r.redis.Incr(ctx, key)
}

func (r *RedisBrimo) Expired(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return r.redis.Expire(ctx, key, expiration)
}
