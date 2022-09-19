package router_test

import (
	"CloudPlatform/pkg/handle"
	"CloudPlatform/pkg/router"
	"CloudPlatform/util"
	"fmt"
	"net/http"

	"testing"
)

var r http.Handler

func init() {
	r = router.Listen()
}

func TestCreateContainer(t *testing.T) {
	param := make(map[string]interface{})
	param["--port"] = "18888:18888"
	param["--name"] = "hello-world"
	param["--image"] = "hello-world:latest"

	req, _ := util.MakeRequest("POST", "JSON", "api/containers", param)
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

func TestGetContainer(t *testing.T) {
	req, _ := util.MakeRequest("GET", "", "api/containers/e5a8c2c553b2", nil)
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

func TestGetContainers(t *testing.T) {

}

func TestRemoveContainer(t *testing.T) {
	req, _ := util.MakeRequest("DELETE", "JSON", "api/containers/e5a8c2c553b2", nil)
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

func TestStartContainer(t *testing.T) {
	req, _ := util.MakeRequest("GET", "", "api/containers/e5a8c2c553b2/start", nil)
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

func TestStopContainer(t *testing.T) {
	req, _ := util.MakeRequest("GET", "", "api/containers/e5a8c2c553b2/stop", nil)
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

func TestRestartContainer(t *testing.T) {
	req, _ := util.MakeRequest("GET", "", "api/containers/e5a8c2c553b2/restart", nil)
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

func TestConnectContainer(t *testing.T) {

}

func TestMakeImage(t *testing.T) {

}

func TestUploadToContainer(t *testing.T) {

}
