package util

import (
	"encoding/json"

	"github.com/go-redis/redis"
	"github.com/golang/glog"
)

var rdb *redis.Client

/** 临时redis, 后期删除 */
func init() {
	rdb = createClient()
}

func GetRDBClient() *redis.Client {
	if rdb == nil {
		rdb = createClient()
	}
	return rdb
}

func createClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "137248",
		DB:       0,
	})

	// 通过 cient.Ping() 来检查是否成功连接到了 redis 服务器
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}

	return client
}

func Hmset(key string, model interface{}) *redis.StatusCmd {
	m := make(map[string]interface{})
	buf, err := json.Marshal(model)
	if err != nil {
		glog.Errorf("redis json Marshal error! msg: %s", err.Error())
		return nil
	}

	if err = json.Unmarshal(buf, &m); err != nil {
		glog.Errorf("redis json Unmarshal error! msg: %s", err.Error())
		return nil
	}
	return rdb.HMSet(key, m)
}

func Hmget(key string) interface{} {
	v := rdb.HMGet(key)
	return v
}

func Get(key string) interface{} {
	result, err := rdb.Get(key).Result()
	if err != nil {
		glog.Errorf("redis get data error! msg: %s", err.Error())
		return nil
	}
	return result
}
