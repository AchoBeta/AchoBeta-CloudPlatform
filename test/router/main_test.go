package router_test

import (
	"cloud-platform/global"

	"cloud-platform/pkg/base"
	"cloud-platform/pkg/base/config"
	commonx "cloud-platform/pkg/handle/common"
	"cloud-platform/pkg/load"
	"cloud-platform/pkg/router"

	_ "cloud-platform/pkg/router/api"
	"fmt"
	"net/http"

	"time"
)

var r http.Handler

const (
	token   = "123456"
	captcha = "1234"
)

func init() {
	load.Init()
	router.RunServer()
	setTokenToRedis()
}

func setTokenToRedis() {
	user := &base.User{
		Id:         "Bb1DaLIAIAA=",
		Username:   "admin",
		Password:   "123456",
		Name:       "marin",
		Containers: []string{containerId},
		Pow:        config.ADMIN_POW,
	}
	str, _ := commonx.StuctToJson(user)
	global.Rdb.Set(fmt.Sprintf(base.TOKEN, token), str, 5*time.Minute)
}

func setCaptchaToRedis() {
	global.Rdb.Set(fmt.Sprintf(base.CAPTCHA, captcha), 1, 5*time.Minute)
}
