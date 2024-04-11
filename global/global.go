package global

import (
	"cloud-platform/pkg/base/cloud"
	"cloud-platform/pkg/base/config"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Config  *config.Server
	Machine *cloud.Machine
	Mgo     *mongo.Client
	Rdb     *redis.Client
)

func GetMgoDb(db string) *mongo.Database {
	return Mgo.Database(db)
}
