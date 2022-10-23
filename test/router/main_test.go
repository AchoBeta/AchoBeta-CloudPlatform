package router_test

import (
	"CloudPlatform/cmd/logic"
	"CloudPlatform/config"
	"CloudPlatform/global"
	"CloudPlatform/internal/base"
	"CloudPlatform/internal/router"
	_ "CloudPlatform/internal/router/api"
	commonx "CloudPlatform/pkg/common"
	"context"
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
	logic.Init("./test_config.yaml")
	r = router.Listen()
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
	global.Rdb.Set(context.TODO(), fmt.Sprintf(base.TOKEN, token), str, 5*time.Minute)
}

func setCaptchaToRedis() {
	global.Rdb.Set(context.TODO(), fmt.Sprintf(base.CAPTCHA, captcha), 1, 5*time.Minute)
}
