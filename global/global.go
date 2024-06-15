package global

import (
	"cloud-platform/internal/base/cloud"
	"cloud-platform/internal/base/config"
	"net/http"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Config     *config.Server
	Machine    *cloud.Machine
	Mgo        *mongo.Client
	Rdb        *redis.Client
	HttpClient *http.Client
)

const (
	LARK_LOGIN_PAGE_URL   = "https://passport.feishu.cn/suite/passport/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code&state=%s"
	LARK_ACCESS_TOKEN_URL = "https://passport.feishu.cn/suite/passport/oauth/token"
	LARK_USERINFO_URL     = "https://passport.feishu.cn/suite/passport/oauth/userinfo"
)

func GetMgoDb(db string) *mongo.Database {
	return Mgo.Database(db)
}
