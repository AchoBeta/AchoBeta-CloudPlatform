package global

import (
	"CloudPlatform/config"
	"CloudPlatform/internal/base"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Config  *config.Server
	Machine *base.Machine
	Mgo     *mongo.Client
	Rdb     *redis.Client
)

func GetMgoDb(db string) *mongo.Database {
	return Mgo.Database(db)
}
