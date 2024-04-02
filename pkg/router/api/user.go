package api

import (
	"cloud-platform/global"
	"cloud-platform/pkg/base"
	"cloud-platform/pkg/handle"
	"cloud-platform/pkg/router/manager"
	"cloud-platform/pkg/service"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/route"
)

func init() {
	manager.RouteHandler.RegisterRouter(manager.LEVEL_GLOBAL, func(router *route.RouterGroup) {
		router.POST("/register", register)
		router.POST("/login", login)
		// todo router.GET("/captcha", captcha1)
	})
	manager.RouteHandler.RegisterRouter(manager.LEVEL_V1, func(router *route.RouterGroup) {
		router.GET("/logout", logout)
	})

	manager.RouteHandler.RegisterRouter(manager.LEVEL_V2, func(router *route.RouterGroup) {
		router.GET("/users", getUsers)
	})
}

func logout(ctx context.Context, c *app.RequestContext) {
	token := c.GetHeader("Authorization")
	r := handle.NewResponse(c)
	_, err := global.Rdb.Del(string(token)).Result()
	if err != nil {
		r.Error(handle.INTERNAL_ERROR)
	} else {
		r.Success(nil)
	}
}

func login(ctx context.Context, c *app.RequestContext) {
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
		hlog.Errorf("[db] del captcha error ! msg: %v\n", err.Error())
		r.Error(handle.INTERNAL_ERROR)
	} else if code == 3 {
		r.Error(handle.USER_ACCOUNT_NOT_EXIST)
	} else if code == 4 {
		hlog.Errorf("user-database decode error! msg: %s", err.Error())
		r.Error(handle.INTERNAL_ERROR)
	} else if code == 5 {
		r.Error(handle.USER_CREDENTIALS_ERROR)
	} else if code == 6 {
		hlog.Errorf("[db] set token to redis error ! msg: %s\n", err.Error())
		r.Error(handle.INTERNAL_ERROR)
	}
}

// Register 注册账号, 成功后返回主键id
func register(ctx context.Context, c *app.RequestContext) {
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
	code, err := service.Register(username, name, password, againPassword, captcha)
	if code == 0 {
		r.Success(nil)
	} else if code == 1 {
		r.Error(handle.CAPTCHA_ERROR)
	} else if code == 2 {
		hlog.Errorf("[db] del captcha error ! msg: %v\n", err.Error())
		r.Error(handle.INTERNAL_ERROR)
	} else if code == 3 {
		r.Error(handle.USER_ACCOUNT_ALREADY_EXIST)
	} else if code == 4 {
		hlog.Errorf("insert user to db error! msg: %s", err.Error())
		r.Error(handle.INTERNAL_ERROR)
	}
}

func getUsers(ctx context.Context, c *app.RequestContext) {
	r := handle.NewResponse(c)
	user := &base.DTOUser{}
	c.BindJSON(user)
	code, users, err := service.GetUsers(user)
	if code == 0 {
		r.Success(users)
	} else if code == 1 {
		hlog.Errorf("[db] find users error ! msg: %s\n", err.Error())
		r.Error(handle.INTERNAL_ERROR)
	} else if code == 2 {
		hlog.Errorf("decode user error ! msg: %s\n", err.Error())
		r.Error(handle.INTERNAL_ERROR)
	}
}

// // 验证码
// func captcha1(ctx context.Context, c *app.RequestContext) {
// 	w, h := 77, 36
// 	captchaId := captcha.NewLen(4)
// 	global.Rdb.Set(fmt.Sprintf(base.CAPTCHA, captchaId), 1, 30*time.Minute)

// 	err := writeResponse(httptest.NewRecorder(), c.Request, captchaId, ".png", "zh", false, w, h)
// 	if err != nil {
// 		hlog.Errorf("create captcha error ! msg: %v\n", err.Error())
// 		r := handle.NewResponse(c)
// 		r.Error(handle.INTERNAL_ERROR)
// 	}
// }

// func writeResponse(w http.ResponseWriter, r *http.Request, id, ext, lang string, download bool, width, height int) error {
// 	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
// 	w.Header().Set("Pragma", "no-cache")
// 	w.Header().Set("Expires", "0")

// 	var content bytes.Buffer
// 	switch ext {
// 	case ".png":
// 		w.Header().Set("Content-Type", "image/png")
// 		_ = captcha.WriteImage(&content, id, width, height)
// 	case ".wav":
// 		w.Header().Set("Content-Type", "audio/x-wav")
// 		_ = captcha.WriteAudio(&content, id, lang)
// 	default:
// 		return captcha.ErrNotFound
// 	}

// 	if download {
// 		w.Header().Set("Content-Type", "application/octet-stream")
// 	}
// 	http.ServeContent(w, r, id+ext, time.Time{}, bytes.NewReader(content.Bytes()))
// 	return nil
// }
