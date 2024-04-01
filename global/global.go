package global

import (
	"cloud-platform/pkg/base/cloud"
	"cloud-platform/pkg/base/config"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var (
	Logger  *zap.SugaredLogger
	Config  *config.Server
	Machine *cloud.Machine
	Mgo     *mongo.Client
	Rdb     *redis.Client
)

func GetMgoDb(db string) *mongo.Database {
	return Mgo.Database(db)
}
