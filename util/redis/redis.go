package redis

import (
	"CloudPlatform/util"
	"context"
	"fmt"
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

// SetStructToHash /**  存入 struct  **/
// HMSet不符合原生命令需要更名
// hmset用于存储结构体
// 封装的这层不要使用日志，通过抛出 error，让调用方决定是否写日志
func SetStructToHash(ctx context.Context, key string, value interface{}, ttl int) ([]redis.Cmder, error) {
	if ttl < 0 {
		return nil, fmt.Errorf("参数传入错误")
	}
	pipe := Rdb.Pipeline()
	pipe.HMSet(ctx, key, util.StructToMap(value))
	// go-redis 用 0 来表示，所以我们也跟随他
	if ttl != 0 {
		pipe.Expire(ctx, key, time.Duration(ttl)*time.Second)
	}
	return pipe.Exec(ctx)
}

// GetHashToStruct /**  获取 struct  **/
// 传入目标target，将get到的属性与值赋予target
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
