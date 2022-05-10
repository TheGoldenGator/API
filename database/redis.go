package database

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Turns a int into a boolean.
// Mostly used for reading boolean values from the database.
func ItoB(i int) bool {
	return i == 1
}

// Check if Redis key expired or not
func CheckCache(ctx context.Context, key string) (bool, error) {
	exists := RDB.Exists(ctx, key)
	return ItoB(int(exists.Val())), nil
}

// Gets cahce of a given key
func GetCache(ctx context.Context, key string) (string, bool) {
	val, err := RDB.Get(ctx, key).Result()
	if err == nil {
		return val, true
	}

	if err != redis.Nil {
		fmt.Println("[INFO] Couldn't query redis")
	}

	return "", false
}

// Sets the cache
func SetCache(ctx context.Context, key, data string, expiration time.Duration) error {
	err := RDB.Set(ctx, key, data, expiration).Err()
	if err != nil {
		return err
	}

	fmt.Println(err)
	return nil
}
