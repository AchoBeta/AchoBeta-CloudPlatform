package test

import (
	"CloudPlatform/base"
	"CloudPlatform/util"
	"fmt"
	"strconv"
)

func Test() {
	u := &base.User{
		Id:       5,
		Username: "avc",
		Password: "sss",
	}
	k := "abcp_user_id_" + strconv.FormatInt(u.Id, 10)
	result := util.Hmset(k, u)
	fmt.Println(result)
}
