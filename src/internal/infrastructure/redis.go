package infrastructure

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

type RedisDB struct {
	Conn *redis.Client
}

var ctx = context.Background()

func NewRedisDB() *RedisDB {
	address := fmt.Sprintf("%s:%s",
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PORT"))

	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &RedisDB{Conn: rdb}
}

func (r *RedisDB) Get(key string) (string, error) {
	val, err := r.Conn.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}

func (r *RedisDB) Set(key string, value string) error {
	err := r.Conn.Set(ctx, key, value, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
