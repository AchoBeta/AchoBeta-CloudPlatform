package global

import (
	"cloud-platform/internal/base/cloud"
	"cloud-platform/internal/base/config"

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
