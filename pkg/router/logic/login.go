package router

import (
<<<<<<< HEAD
	"CloudPlatform/base"
	"CloudPlatform/conf/secret"
	"CloudPlatform/pkg/handle"
	"CloudPlatform/pkg/router"
	"CloudPlatform/util"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

const ABCP_USER_KEY string = "abcp_user_id_"

func init() {
	router.Register(func(router gin.IRoutes) {
		router.POST("/register", Register)
		router.POST("/login", Login)
=======
	"CloudPlatform/pkg/handle"
	"CloudPlatform/pkg/router"

	"github.com/gin-gonic/gin"
)

func init() {
	router.Register(func(router gin.IRoutes) {
		router.GET("/test", test)
>>>>>>> master
	}, router.V0)

	router.Register(func(router gin.IRoutes) {
		router.GET("/test")
	}, router.V1)
}

<<<<<<< HEAD
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	r := handle.NewResponse(c)

	if password == "" || username == "" {
		r.Error(handle.PARAM_IS_BLANK)
		return
	}

	// 通过username获取uid 再获取 user 的所有信息
	byte, err := json.Marshal(util.Hmget(util.Get(username).(string)))
	if err != nil {
		glog.Errorf("interface json marshal error! msg: %s", err.Error())
		r.Error(handle.INTERNAL_ERROR)
		return
	}
	var user base.User
	if jsonRes := json.Unmarshal(byte, &user); jsonRes != nil {
		glog.Errorf("user json Unmarshal error! msg: %s", err.Error())
		r.Error(handle.INTERNAL_ERROR)
		return
	}

	if password != string(secret.Decrypt(user.Password)) {
		r.Error(handle.USER_CREDENTIALS_ERROR)
		return
	}
	r.Success(createToken(user.Id))
}

// Register 注册账号, 成功后返回主键id
func Register(c *gin.Context) {
	rdb := util.GetRDBClient()
	username := c.PostForm("username")
	password := c.PostForm("password")
	againPassword := c.PostForm("againPassword")
	r := handle.NewResponse(c)
	if password != againPassword {
		r.Error(handle.USER_PASSWORD_DIFFERENT)
		return
	}

	// 通过redis生成id, 保证全局唯一自增id
	id, err := rdb.Incr("user_id_incr_").Result()
	if err != nil {
		glog.Errorf("redis incr error, msg: %s", err.Error())
		return
	}
	// 这里先做一个账号数据插入, 具体信息后续再设置
	userModel := &base.User{
		Id:       id,
		Username: username,
		Password: secret.Encrypt(password),
	}
	// 插入数据库, redis做数据库，需要额外存一个username - uid 的数据
	key := ABCP_USER_KEY + strconv.FormatInt(id, 10)

	rdb.Set(username, key, -1).Result()
	result, err := util.Hmset(key, userModel).Result()

	if err != nil {
		glog.Errorf("insert user info error!")
		r.Error(handle.COMMON_FAIL)
		return
	}

	r.Success(result)
}

/** 私有方法 */
func createToken(id int64) string {
	snowId := util.CreateSnowflake(id)
	if snowId == -1 {
		return ""
	}
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(snowId))
	return base64.StdEncoding.EncodeToString(buf)
=======
func test(c *gin.Context) {
	r := handle.NewResponse(c)
	r.Success(nil)
>>>>>>> master
}
