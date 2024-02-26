package cache

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var (
	Ctx    = context.Background()
	Client = NewClient()
)

func NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func Set(key string, value int64) error {
	err := Client.Set(Ctx, key, value, 0).Err()
	return err
}

func Get(key string) (int64, error) {
	val, err := Client.Get(Ctx, key).Result()
	if err != nil {
		return 0, err
	}

	intVal, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, err
	}

	return intVal, nil
}
