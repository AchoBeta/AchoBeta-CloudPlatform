package router_test

import (
	"CloudPlatform/config"
	"CloudPlatform/internal/base"
	"CloudPlatform/internal/handle"
	requestx "CloudPlatform/pkg/request"
	"testing"
)

func TestRegistry(t *testing.T) {
	form := make(map[string]string)
	form["username"] = "admin"
	form["password"] = "123456"
	form["againPassword"] = "123456"
	form["name"] = "marin"
	form["captcha"] = captcha
	setCaptchaToRedis()
	req, _ := requestx.MakeRequest("POST", "form", "/api/register", form)
	resp := handle.JsonMsgResult{}
	err := requestx.Request(r, req, &resp)
	if err != nil {
		t.Error(err)
	} else if resp.Code != 200 {
		t.Error(resp.Message)
	} else {
		t.Log(resp)
	}
}

func TestLogin(t *testing.T) {
	form := make(map[string]string)
	form["username"] = "admin"
	form["password"] = "123456"
	form["captcha"] = captcha
	setCaptchaToRedis()
	req, _ := requestx.MakeRequest("POST", "form", "/api/login", form)
	resp := handle.JsonMsgResult{}
	err := requestx.Request(r, req, &resp)
	if err != nil {
		t.Error(err)
	} else if resp.Code != 200 {
		t.Error(resp.Message)
	} else {
		t.Log(resp)
	}
}

func TestLogout(t *testing.T) {
	req, _ := requestx.MakeRequest("GET", "json", "/api/logout", nil)
	req.Header.Set("Authorization", token)
	resp := handle.JsonMsgResult{}
	err := requestx.Request(r, req, &resp)
	if err != nil {
		t.Error(err)
	} else if resp.Code != 200 {
		t.Error(resp.Message)
	} else {
		t.Log(resp)
	}
}

func TestGetUsers(t *testing.T) {
	user := base.DTOUser{
		Pow: config.TOURIST_POW,
	}
	req, _ := requestx.MakeRequest("GET", "json", "/api/users", &user)
	req.Header.Set("Authorization", token)
	resp := handle.JsonMsgResult{}
	err := requestx.Request(r, req, &resp)
	if err != nil {
		t.Error(err)
	} else if resp.Code != 200 {
		t.Error(resp.Message)
	} else {
		t.Log(resp)
	}
}
