package api

import (
	"CloudPlatform/global"
	"CloudPlatform/internal/base"
	"CloudPlatform/internal/handle"
	"CloudPlatform/internal/router"
	"CloudPlatform/internal/service"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func init() {
	router.Register(func(router gin.IRoutes) {
		router.POST("/register", register)
		router.POST("/login", login)
		router.GET("/captcha", captcha1)
	}, router.V0)

	router.Register(func(router gin.IRoutes) {
		router.GET("/logout", logout)
	}, router.V1)

	router.Register(func(router gin.IRoutes) {
		router.GET("/users", getUsers)
	}, router.V2)
}

func logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	r := handle.NewResponse(c)
	_, err := global.Rdb.Del(context.TODO(), token).Result()
	if err != nil {
		r.Error(handle.INTERNAL_ERROR)
	} else {
		r.Success(nil)
	}
}

func login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	captcha := c.PostForm("captcha")
	r := handle.NewResponse(c)
	if password == "" || username == "" || captcha == "" {
		r.Error(handle.PARAM_IS_BLANK)
		return
	}
	user := &base.DTOUser{}
	err, code, token := service.Login(username, password, captcha, user)
	if code == 0 {
		r.Ctx.Header("Authorization", token)
		r.Success(user)
	} else if code == 1 {
		r.Error(handle.CAPTCHA_ERROR)
	} else if code == 2 {
		glog.Errorf("[db] del captcha error ! msg: %v\n", err.Error())
		r.Error(handle.INTERNAL_ERROR)
	} else if code == 3 {
		r.Error(handle.USER_ACCOUNT_NOT_EXIST)
	} else if code == 4 {
		glog.Errorf("user-database decode error! msg: %s", err.Error())
		r.Error(handle.INTERNAL_ERROR)
	} else if code == 5 {
		r.Error(handle.USER_CREDENTIALS_ERROR)
	} else if code == 6 {
		glog.Errorf("[db] set token to redis error ! msg: %s\n", err.Error())
		r.Error(handle.INTERNAL_ERROR)
	}
}

// Register 注册账号, 成功后返回主键id
func register(c *gin.Context) {
	username := c.PostForm("username")
	name := c.PostForm("name")
	password := c.PostForm("password")
	againPassword := c.PostForm("againPassword")
	captcha := c.PostForm("captcha")
	r := handle.NewResponse(c)
	if username == "" || password == "" || name == "" || captcha == "" {
		r.Error(handle.PARAM_NOT_COMPLETE)
		return
	}
	if password != againPassword {
		r.Error(handle.USER_PASSWORD_DIFFERENT)
		return
	}
	err, code := service.Register(username, name, password, againPassword, captcha)
	if code == 0 {
		r.Success(nil)
	} else if code == 1 {
		r.Error(handle.CAPTCHA_ERROR)
	} else if code == 2 {
		glog.Errorf("[db] del captcha error ! msg: %v\n", err.Error())
		r.Error(handle.INTERNAL_ERROR)
	} else if code == 3 {
		r.Error(handle.USER_ACCOUNT_ALREADY_EXIST)
	} else if code == 4 {
		glog.Errorf("insert user to db error! msg: %s", err.Error())
		r.Error(handle.INTERNAL_ERROR)
	}
}

// 验证码
func captcha1(c *gin.Context) {
	w, h := 77, 36
	captchaId := captcha.NewLen(4)
	global.Rdb.Set(context.TODO(), fmt.Sprintf(base.CAPTCHA, captchaId), 1, 30*time.Minute)
	err := writeResponse(c.Writer, c.Request, captchaId, ".png", "zh", false, w, h)
	if err != nil {
		glog.Errorf("create captcha error ! msg: %v\n", err.Error())
		r := handle.NewResponse(c)
		r.Error(handle.INTERNAL_ERROR)
	}
}

func getUsers(c *gin.Context) {
	r := handle.NewResponse(c)
	user := &base.DTOUser{}
	c.BindJSON(user)
	err, code, users := service.GetUsers(user)
	if code == 0 {
		r.Success(users)
	} else if code == 1 {
		glog.Errorf("[db] find users error ! msg: %s\n", err.Error())
		r.Error(handle.INTERNAL_ERROR)
	} else if code == 2 {
		glog.Errorf("decode user error ! msg: %s\n", err.Error())
		r.Error(handle.INTERNAL_ERROR)
	}
}

func writeResponse(w http.ResponseWriter, r *http.Request, id, ext, lang string, download bool, width, height int) error {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	var content bytes.Buffer
	switch ext {
	case ".png":
		w.Header().Set("Content-Type", "image/png")
		_ = captcha.WriteImage(&content, id, width, height)
	case ".wav":
		w.Header().Set("Content-Type", "audio/x-wav")
		_ = captcha.WriteAudio(&content, id, lang)
	default:
		return captcha.ErrNotFound
	}

	if download {
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	http.ServeContent(w, r, id+ext, time.Time{}, bytes.NewReader(content.Bytes()))
	return nil
}
