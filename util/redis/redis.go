package redis

import (
	"CloudPlatform/util"
	"context"
	"time"

	"github.com/mitchellh/mapstructure"

	"github.com/go-redis/redis/v9"
)

// Rdb Background返回一个非空的Context。 它永远不会被取消，没有值，也没有期限。
// 它通常在main函数，初始化和测试时使用，并用作传入请求的顶级上下文。
var Rdb *redis.Client

// Connect /**      连接客户端        */
func Connect(ctx context.Context) error {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     RDb,
		Password: RDbPwd,
		DB:       RDbNumber,
	})

	_, err := Rdb.Ping(ctx).Result()
	return err
}

// SetStructToHash
// Zero expiration means the key has no expiration time.
func SetStructToHash(ctx context.Context, key string, value interface{}, expiration time.Duration) ([]redis.Cmder, error) {
	pipe := Rdb.Pipeline()
	pipe.HMSet(ctx, key, util.StructToMap(value))
	pipe.Expire(ctx, key, expiration)
	return pipe.Exec(ctx)
}

// GetHashToStruct
func GetHashToStruct(ctx context.Context, key string, target interface{}) error {
	result, err := Rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return err
	}
	err = mapstructure.Decode(result, &target)
	if err != nil {
		return err
	}
	return nil
}
