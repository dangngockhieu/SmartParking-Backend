package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"backend/configs"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
	Env    string
}

func NewRedis(cfg *configs.Config) *RedisClient {
	ctx := context.Background()

	addr := cfg.RedisAddr
	if addr == "" {
		addr = "127.0.0.1:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Cannot connect to redis: %v", err)
	}

	return &RedisClient{
		Client: client,
		Ctx:    ctx,
		Env:    cfg.AppEnv,
	}
}

func (r *RedisClient) Close() error {
	return r.Client.Close()
}

func (r *RedisClient) Key(parts ...any) string {
	key := r.Env
	for _, part := range parts {
		key += ":" + fmt.Sprint(part)
	}
	return key
}

func (r *RedisClient) Set(key string, value any, ttl time.Duration) error {
	return r.Client.Set(r.Ctx, key, value, ttl).Err()
}

func (r *RedisClient) SetNX(key string, value any, ttl time.Duration) (bool, error) {
	return r.Client.SetNX(r.Ctx, key, value, ttl).Result()
}

func (r *RedisClient) Get(key string) (string, error) {
	return r.Client.Get(r.Ctx, key).Result()
}

func (r *RedisClient) Delete(key string) error {
	return r.Client.Del(r.Ctx, key).Err()
}

func (r *RedisClient) Exists(key string) (bool, error) {
	n, err := r.Client.Exists(r.Ctx, key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *RedisClient) TTL(key string) (time.Duration, error) {
	return r.Client.TTL(r.Ctx, key).Result()
}

func (r *RedisClient) Expire(key string, ttl time.Duration) error {
	return r.Client.Expire(r.Ctx, key, ttl).Err()
}

func (r *RedisClient) Incr(key string) (int64, error) {
	return r.Client.Incr(r.Ctx, key).Result()
}

func (r *RedisClient) Decr(key string) (int64, error) {
	return r.Client.Decr(r.Ctx, key).Result()
}

func (r *RedisClient) HSet(key string, values ...any) error {
	return r.Client.HSet(r.Ctx, key, values...).Err()
}

func (r *RedisClient) HGet(key, field string) (string, error) {
	return r.Client.HGet(r.Ctx, key, field).Result()
}

func (r *RedisClient) HGetAll(key string) (map[string]string, error) {
	return r.Client.HGetAll(r.Ctx, key).Result()
}

func (r *RedisClient) HIncrBy(key, field string, incr int64) (int64, error) {
	return r.Client.HIncrBy(r.Ctx, key, field, incr).Result()
}
