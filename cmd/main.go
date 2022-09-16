package main

import (
	"CloudPlatform/base"
	"CloudPlatform/cmd/logic"
	"CloudPlatform/util"
	"fmt"
	"strconv"
)

func main() {
	logic.Run()
}

func test() {
	u := &base.User{
		Id:       5,
		Username: "avc",
		Password: "sss",
	}
	k := "abcp_user_id_" + strconv.FormatInt(u.Id, 10)
	result := util.Hmset(k, u)
	fmt.Println(result)
}
