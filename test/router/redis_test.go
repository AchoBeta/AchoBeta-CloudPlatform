package router_test

import (
	"cloud-platform/global"
	redisx "cloud-platform/pkg/handle/redis"
	"context"
	"fmt"
	"testing"
)

type Student struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

func TestHSet(t *testing.T) {
	var ctx = context.Background()
	defer global.Rdb.Close()
	stu := Student{
		Name: "1231",
		Id:   "1231231",
	}
	cmd, err := redisx.SetStructToRedis(global.Rdb, "test920", stu, 0)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(cmd)
	//Del("test920:")
	var stud Student
	redisx.GetStruceFromRedis(ctx, global.Rdb, "test920", &stud)
	fmt.Println(stud)
	redisx.GetStruceFromRedis(ctx, global.Rdb, "123", "123")
}
