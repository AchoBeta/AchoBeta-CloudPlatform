package router_test

import (
	"CloudPlatform/pkg/handle"
	"CloudPlatform/util"
	"fmt"
	"testing"
)

func TestGetImages(t *testing.T) {
	req, _ := util.MakeRequest("GET", "", "/images", nil)
	resp := handle.JsonMsgResult{}

	err := util.Request(r, req, &resp)
	if err != nil {
		t.Error(err)
		return
	} else if resp.Code != 200 {
		t.Errorf(resp.Message)
		return
	} else {
		fmt.Println(resp)
	}
}

func TestGetImage(t *testing.T) {
	req, _ := util.MakeRequest("GET", "", "/images/d7a966a74f16", nil)
	resp := handle.JsonMsgResult{}

	err := util.Request(r, req, &resp)
	if err != nil {
		t.Error(err)
		return
	} else if resp.Code != 200 {
		t.Errorf(resp.Message)
		return
	} else {
		fmt.Println(resp)
	}
}

func TestBuildImage(t *testing.T) {

}

func TestDeleteImage(t *testing.T) {
	req, _ := util.MakeRequest("DELETE", "", "/images/d7a966a74f16", nil)
	resp := handle.JsonMsgResult{}

	err := util.Request(r, req, &resp)
	if err != nil {
		t.Error(err)
		return
	} else if resp.Code != 200 {
		t.Errorf(resp.Message)
		return
	} else {
		fmt.Println(resp)
	}
}

func TestSearchImages(t *testing.T) {
	param := make(map[string]interface{})
	param["image"] = "hello-world"
	param["tag"] = "latest"

	req, _ := util.MakeRequest("GET", "", "/images/search", param)
	resp := handle.JsonMsgResult{}

	err := util.Request(r, req, &resp)
	if err != nil {
		t.Error(err)
		return
	} else if resp.Code != 200 {
		t.Errorf(resp.Message)
		return
	} else {
		fmt.Println(resp)
	}
}

func TestPushImage(t *testing.T) {
	param := make(map[string]interface{})
	param["image"] = "hello-world"
	param["tag"] = "latest"

	req, _ := util.MakeRequest("GET", "", "/images/push", param)
	resp := handle.JsonMsgResult{}

	err := util.Request(r, req, &resp)
	if err != nil {
		t.Error(err)
		return
	} else if resp.Code != 200 {
		t.Errorf(resp.Message)
		return
	} else {
		fmt.Println(resp)
	}
}
