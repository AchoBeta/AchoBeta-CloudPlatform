package exec

import (
	"cloud-platform/global"
	"context"

	"github.com/golang/glog"
)

func Eve() {
	glog.Flush()
	global.Rdb.Close()
	global.Mgo.Disconnect(context.TODO())
}
