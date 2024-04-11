package router_test

import (
	"cloud-platform/pkg/base/cloud"
	"cloud-platform/pkg/handle"
	requestx "cloud-platform/pkg/handle/request"
	"testing"
)

const containerId = "5ffcdf0c90fec7f75434b0d8084241577c38967a52af76362dfbd7c24124624d"

func TestCreateContainer(t *testing.T) {
	container := cloud.Container{
		Name:  "base",
		Image: "achobeta/abcp_base:0.1",
		Param: cloud.Param{
			Ports: []int{18888},
		},
	}
	req, _ := requestx.MakeRequest("POST", "JSON", "/api/containers", container)
	req.Header.Set("Authorization", token)
	resp := handle.JsonMsgResult{}

	err := requestx.Request(r, req, &resp)
	if err != nil {
		t.Error(err)
	} else if resp.Code != 200 {
		t.Errorf(resp.Message)
	} else {
		t.Log(resp)
	}
}

func TestGetContainer(t *testing.T) {
	req, _ := requestx.MakeRequest("GET", "json", "/api/containers/"+containerId, nil)
	req.Header.Set("Authorization", token)
	resp := handle.JsonMsgResult{}

	err := requestx.Request(r, req, &resp)
	if err != nil {
		t.Error(err)
	} else if resp.Code != 200 {
		t.Errorf(resp.Message)
	} else {
		t.Log(resp)
	}
}

func TestGetContainers(t *testing.T) {
	req, _ := requestx.MakeRequest("GET", "json", "/api/containers", nil)
	req.Header.Set("Authorization", token)
	resp := handle.JsonMsgResult{}

	err := requestx.Request(r, req, &resp)
	if err != nil {
		t.Error(err)
	} else if resp.Code != 200 {
		t.Errorf(resp.Message)
	} else {
		t.Log(resp)
	}
}

func TestStartContainer(t *testing.T) {
	req, _ := requestx.MakeRequest("GET", "json", "/api/containers/"+containerId+"/start", nil)
	req.Header.Set("Authorization", token)
	resp := handle.JsonMsgResult{}

	err := requestx.Request(r, req, &resp)
	if err != nil {
		t.Error(err)
	} else if resp.Code != 200 {
		t.Errorf(resp.Message)
	} else {
		t.Log(resp)
	}
}

func TestStopContainer(t *testing.T) {
	req, _ := requestx.MakeRequest("GET", "json", "/api/containers/"+containerId+"/stop", nil)
	req.Header.Set("Authorization", token)
	resp := handle.JsonMsgResult{}

	err := requestx.Request(r, req, &resp)
	if err != nil {
		t.Error(err)
	} else if resp.Code != 200 {
		t.Errorf(resp.Message)
	} else {
		t.Log(resp)
	}
}

func TestRestartContainer(t *testing.T) {
	req, _ := requestx.MakeRequest("GET", "", "api/containers/"+containerId+"/restart", nil)
	resp := handle.JsonMsgResult{}

	err := requestx.Request(r, req, &resp)
	if err != nil {
		t.Error(err)
	} else if resp.Code != 200 {
		t.Errorf(resp.Message)
	} else {
		t.Log(resp)
	}
}

func TestRemoveContainer(t *testing.T) {
	req, _ := requestx.MakeRequest("DELETE", "JSON", "/api/containers/"+containerId, nil)
	req.Header.Set("Authorization", token)
	resp := handle.JsonMsgResult{}

	err := requestx.Request(r, req, &resp)
	if err != nil {
		t.Error(err)
	} else if resp.Code != 200 {
		t.Errorf(resp.Message)
	} else {
		t.Log(resp)
	}
}

func TestMakeImage(t *testing.T) {
	image := cloud.Image{
		Name:   "abcp/base1",
		Tag:    "test",
		Author: "marin",
		Desc:   "test",
	}
	req, _ := requestx.MakeRequest("POST", "JSON", "/api/containers/"+containerId+"/makeImage", image)
	req.Header.Set("Authorization", token)
	resp := handle.JsonMsgResult{}

	err := requestx.Request(r, req, &resp)
	if err != nil {
		t.Error(err)
	} else if resp.Code != 200 {
		t.Errorf(resp.Message)
	} else {
		t.Log(resp)
	}
}
