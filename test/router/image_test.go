package router_test

import (
	"CloudPlatform/internal/handle"
	requestx "CloudPlatform/pkg/request"
	"fmt"
	"testing"
)

const imageId = "ad77059867bd"

func TestGetImages(t *testing.T) {
	req, _ := requestx.MakeRequest("GET", "json", "/api/images", nil)
	req.Header.Set("Authorization", token)
	resp := handle.JsonMsgResult{}

	err := requestx.Request(r, req, &resp)
	if err != nil {
		t.Error(err)
	} else if resp.Code != 200 {
		t.Errorf(resp.Message)
	} else {
		fmt.Println(resp)
	}
}

func TestGetImage(t *testing.T) {
	req, _ := requestx.MakeRequest("GET", "json", "/api/images/"+imageId, nil)
	req.Header.Set("Authorization", token)
	resp := handle.JsonMsgResult{}

	err := requestx.Request(r, req, &resp)
	if err != nil {
		t.Error(err)
	} else if resp.Code != 200 {
		t.Errorf(resp.Message)
	} else {
		fmt.Println(resp)
	}
}

func TestBuildImage(t *testing.T) {

}

func TestDeleteImage(t *testing.T) {
	req, _ := requestx.MakeRequest("DELETE", "json", "/api/images/"+imageId, nil)
	req.Header.Set("Authorization", token)
	resp := handle.JsonMsgResult{}

	err := requestx.Request(r, req, &resp)
	if err != nil {
		t.Error(err)
	} else if resp.Code != 200 {
		t.Errorf(resp.Message)
	} else {
		fmt.Println(resp)
	}
}

func TestSearchImages(t *testing.T) {
	param := make(map[string]interface{})
	param["image"] = "hello-world"
	param["tag"] = "latest"

	req, _ := requestx.MakeRequest("GET", "json", "/images/search", param)
	resp := handle.JsonMsgResult{}

	err := requestx.Request(r, req, &resp)
	if err != nil {
		t.Error(err)
	} else if resp.Code != 200 {
		t.Errorf(resp.Message)
	} else {
		fmt.Println(resp)
	}
}

func TestPushImage(t *testing.T) {
	param := make(map[string]interface{})
	param["image"] = "hello-world"
	param["tag"] = "latest"

	req, _ := requestx.MakeRequest("GET", "json", "/images/push", param)
	resp := handle.JsonMsgResult{}

	err := requestx.Request(r, req, &resp)
	if err != nil {
		t.Error(err)
	} else if resp.Code != 200 {
		t.Errorf(resp.Message)
	} else {
		fmt.Println(resp)
	}
}
