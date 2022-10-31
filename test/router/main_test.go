package router_test

import (
	"cloud-platform/global"
	"cloud-platform/internal/base"
	"cloud-platform/internal/base/config"
	"cloud-platform/internal/exec"
	commonx "cloud-platform/internal/pkg/common"
	"cloud-platform/internal/router"
	_ "cloud-platform/internal/router/api"
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
	exec.Init("./test_config.yaml")
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
	str, _ := commonx.StructToJson(user)
	global.Rdb.Set(context.TODO(), fmt.Sprintf(base.TOKEN, token), str, 5*time.Minute)
}

func setCaptchaToRedis() {
	global.Rdb.Set(fmt.Sprintf(base.CAPTCHA, captcha), 1, 5*time.Minute)
}
