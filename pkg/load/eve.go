package load

import (
	"cloud-platform/global"
	"context"
)

func Eve() {
	global.Rdb.Close()
	global.Mgo.Disconnect(context.TODO())
}
