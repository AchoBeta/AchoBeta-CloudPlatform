package redis

import (
	"gopkg.in/ini.v1"
	"log"
)

var (
	RDb       string
	RDbPwd    string
	RDbNumber int
)

func init() {
	file, err := ini.Load("../redis/conf.ini")
	if err != nil {
		log.Fatal(err)
	}
	loadRedisData(file)
}

func loadRedisData(file *ini.File) {
	RDb = file.Section("redis").Key("RedisAddr").String()
	RDbPwd = file.Section("redis").Key("RedisPw").String()
	RDbNumber, _ = file.Section("Redis").Key("RedisDb").Int()
}
