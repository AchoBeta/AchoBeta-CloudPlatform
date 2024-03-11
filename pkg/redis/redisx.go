package redisx

import (
	commonx "CloudPlatform/pkg/common"
	"context"
	"github.com/mitchellh/mapstructure"
	"github.com/redis/go-redis/v9"
	"time"
)

// SetStructToHash
// Zero expiration means the key has no expiration time.
func SetStructToRedis(ctx context.Context, rdb *redis.Client, key string, value interface{}, expiration time.Duration) ([]redis.Cmder, error) {
	pipe := rdb.Pipeline()
	pipe.HMSet(ctx, key, commonx.StructToMap(value))
	pipe.Expire(ctx, key, expiration)
	return pipe.Exec(ctx)
}

// GetHashToStruct
func GetStruceFromRedis(ctx context.Context, rdb *redis.Client, key string, target interface{}) error {
	result, err := rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return err
	}
	err = mapstructure.Decode(result, &target)
	if err != nil {
		return err
	}
	return nil
}
