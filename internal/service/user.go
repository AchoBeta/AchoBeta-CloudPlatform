package service

import (
	"cloud-platform/global"
	"cloud-platform/internal/base"
	"cloud-platform/internal/base/config"
	"cloud-platform/internal/handle"
	commonx "cloud-platform/internal/pkg/common"
	requestx "cloud-platform/internal/pkg/request"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/go-redis/redis"
	"github.com/golang/glog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Login(username, password, captcha string, dtoUser *base.DTOUser) (int8, string, error) {
	// 判断验证码是否正确
	cmd := global.Rdb.Del(fmt.Sprintf(base.CAPTCHA, captcha))
	if cmd.Err() != nil {
		if cmd.Err() == redis.Nil {
			return 1, "", cmd.Err()
		} else {
			return 2, "", cmd.Err()
		}
	}
	// 判断数据库是否有此用户
	filter := bson.M{"username": username}
	res := global.GetMgoDb("abcp").Collection("user").FindOne(context.TODO(), filter)
	if res.Err() != nil {
		return 3, "", res.Err()
	}
	var user base.User
	err := res.Decode(&user)
	if err != nil {
		return 4, "", err
	}
	// 验证密码是否正确
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%s-%s", password, global.Config.App.Salt)))
	if user.Password != fmt.Sprintf("%x", h.Sum(nil)) {
		return 5, "", err
	}
	token := createToken()
	str, _ := commonx.StructToJson(&user)
	cmd1 := global.Rdb.Set(fmt.Sprintf(base.TOKEN, token), str, 30*time.Hour)
	if cmd1.Err() != nil {
		return 6, "", cmd1.Err()
	}
	dtoUser.Id = user.Id
	dtoUser.Username = user.Username
	dtoUser.Name = user.Name
	dtoUser.Pow = "tourist"
	if user.Pow == config.ADMIN_POW {
		dtoUser.Pow = "admin"
	} else if user.Pow == config.USER_POW {
		dtoUser.Pow = "user"
	}
	dtoUser.Containers = user.Containers
	return 0, token, nil
}

func Register(username, name, password, againPassword, captcha string) (int8, error) {
	// 判断验证码是否正确
	cmd := global.Rdb.Del(fmt.Sprintf(base.CAPTCHA, captcha))
	if cmd.Err() != nil {
		if cmd.Err() == redis.Nil {
			return 1, cmd.Err()
		} else {
			return 2, cmd.Err()
		}
	}
	// 判断用户是否存在
	filter := bson.M{"username": username}
	res := global.GetMgoDb("abcp").Collection("user").FindOne(context.TODO(), filter)
	if res.Err() != mongo.ErrNoDocuments {
		return 3, res.Err()
	}
	// 添加数据
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%s-%s", password, global.Config.App.Salt)))
	user := base.User{
		Id:         createToken(),
		Username:   username,
		Name:       name,
		Password:   fmt.Sprintf("%x", h.Sum(nil)),
		Pow:        config.TOURIST_POW,
		Containers: []string{},
	}
	_, err := global.GetMgoDb("abcp").Collection("user").InsertOne(context.TODO(), &user)
	if err != nil {
		return 4, err
	}
	return 0, nil
}

func createToken() string {
	snowId := commonx.GetNextSnowflakeID()
	if snowId == -1 {
		return ""
	}
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(snowId))
	return base64.StdEncoding.EncodeToString(buf)
}

func GetUsers(user *base.DTOUser) (int8, []base.DTOUser, error) {
	collection := global.GetMgoDb("abcp").Collection("user")
	cur, err := collection.Find(context.TODO(), user)
	if err != nil {
		return 1, nil, err
	}
	defer cur.Close(context.TODO())
	users := []base.DTOUser{}
	for cur.Next(context.TODO()) {
		user := base.User{}
		err = cur.Decode(&user)
		if err != nil {
			return 2, nil, err
		}
		pow := "admin"
		if user.Pow == config.USER_POW {
			pow = "user"
		} else if user.Pow == config.TOURIST_POW {
			pow = user.Pow
		}
		dtoUser := base.DTOUser{
			Id:         user.Id,
			Username:   user.Username,
			Name:       user.Name,
			Pow:        pow,
			Containers: user.Containers,
		}
		users = append(users, dtoUser)
	}
	return 0, users, nil
}

func LarkLogin(code string) {
	param := map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     global.Config.App.Lark.AppId,
		"client_secret": global.Config.App.Lark.AppSecret,
		"code":          code,
		"redirect_uri":  global.Config.App.Lark.RedirectUrl,
	}
	req, _ := requestx.MakeRequest("POST", "form", global.LARK_ACCESS_TOKEN_URL, param)
	result, err := global.HttpClient.Do(req)
	if err != nil {
		return
	}
	defer result.Body.Close()

	// extract response body
	resp := handle.JsonMsgResult{}
	bodyByte, _ := io.ReadAll(result.Body)
	json.Unmarshal(bodyByte, &resp)
	glog.Errorf("%v", resp)
}
