package test

import (
	"CloudPlatform/util/redis"
	"context"
	"fmt"
	"testing"

	"github.com/golang/glog"
)

type Student struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

func TestHSet(t *testing.T) {
	var ctx = context.Background()
	err := redis.Connect(ctx)
	if err != nil {
		glog.Error(err)
		return
	}
	fmt.Println(redis.Rdb)
	defer redis.Rdb.Close()
	stu := Student{
		Name: "1231",
		Id:   "1231231",
	}
	cmd, err := redis.SetStructToHash(ctx, "test920", stu, 0)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(cmd)
	//Del("test920:")
	var stud Student
	redis.GetHashToStruct(ctx, "test920", &stud)
	fmt.Println(stud)
	redis.GetHashToStruct(ctx, "123", "123")
}
