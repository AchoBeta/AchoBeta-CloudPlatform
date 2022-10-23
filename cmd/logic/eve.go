package logic

import (
	"CloudPlatform/global"
	"context"

	"github.com/golang/glog"
)

func Eve() {
	glog.Flush()
	global.Rdb.Close()
	global.Mgo.Disconnect(context.TODO())
}
