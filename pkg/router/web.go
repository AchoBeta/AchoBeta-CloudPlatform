package router

import (
	"cloud-platform/global"
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

const (
	V0 uint8 = 0
	V1 uint8 = 1
	V2 uint8 = 2
	V3 uint8 = 3
)

var (
	_hooks_V0, _hooks_V1, _hooks_V2, _hooks_V3 []Hook
)

type Hook func(ctx context.Context, c *app.RequestContext)

func Register(hook Hook, hookType uint8) {
	switch hookType {
	case V0:
		_hooks_V0 = append(_hooks_V0, hook)
	case V1:
		_hooks_V1 = append(_hooks_V1, hook)
	case V2:
		_hooks_V2 = append(_hooks_V2, hook)
	case V3:
		_hooks_V3 = append(_hooks_V3, hook)
	default:
		global.Logger.Error("Register Error")
	}
}

func Run() {
	h, err := Listen()
	if err != nil {
		global.Logger.Errorf("Listen error: %v", err)
		panic(err.Error())
	}
	h.Spin()
}

func Listen() (*server.Hertz, error) {
	h := server.Default(server.WithHostPorts(fmt.Sprintf("%s:%d", global.Config.App.Host, global.Config.App.Port)))

	return h, nil
}
