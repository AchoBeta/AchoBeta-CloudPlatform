package redisx

import (
	"context"
	"time"

	commonx "cloud-platform/internal/pkg/common"

	"github.com/go-redis/redis"
	"github.com/mitchellh/mapstructure"
)

// SetStructToHash
// Zero expiration means the key has no expiration time.
func SetStructToRedis(rdb *redis.Client, key string, value interface{}, expiration time.Duration) ([]redis.Cmder, error) {
	pipe := rdb.Pipeline()
	pipe.HMSet(key, commonx.StructToMap(value))
	pipe.Expire(key, expiration)
	return pipe.Exec()
}

// GetHashToStruct
func GetStruceFromRedis(ctx context.Context, rdb *redis.Client, key string, target interface{}) error {
	result, err := rdb.HGetAll(key).Result()
	if err != nil {
		return err
	}
	err = mapstructure.Decode(result, &target)
	if err != nil {
		return err
	}
	return nil
}
