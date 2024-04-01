package service

import (
	"cloud-platform/global"
	"cloud-platform/pkg/base"
	"cloud-platform/pkg/base/config"
	commonx "cloud-platform/pkg/handle/common"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Login(username, password, captcha string, dtoUser *base.DTOUser) (error, int8, string) {
	// 判断验证码是否正确
	cmd := global.Rdb.Del(fmt.Sprintf(base.CAPTCHA, captcha))
	if cmd.Err() != nil {
		if cmd.Err() == redis.Nil {
			return cmd.Err(), 1, ""
		} else {
			return cmd.Err(), 2, ""
		}
	}
	// 判断数据库是否有此用户
	filter := bson.M{"username": username}
	res := global.GetMgoDb("abcp").Collection("user").FindOne(context.TODO(), filter)
	if res.Err() != nil {
		return res.Err(), 3, ""
	}
	var user base.User
	err := res.Decode(&user)
	if err != nil {
		return err, 4, ""
	}
	// 验证密码是否正确
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%s-%s", password, global.Config.App.Salt)))
	if user.Password != fmt.Sprintf("%x", h.Sum(nil)) {
		return err, 5, ""
	}
	token := createToken()
	str, _ := commonx.StuctToJson(user)
	cmd1 := global.Rdb.Set(fmt.Sprintf(base.TOKEN, token), str, 30*time.Minute)
	if cmd1.Err() != nil {
		return cmd1.Err(), 6, ""
	}
	dtoUser.Id = user.Id
	dtoUser.Username = user.Username
	dtoUser.Name = user.Name
	dtoUser.Pow = user.Pow
	dtoUser.Containers = user.Containers
	return nil, 0, token
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
		dtoUser := base.DTOUser{
			Id:         user.Id,
			Username:   user.Username,
			Name:       user.Name,
			Pow:        user.Pow,
			Containers: user.Containers,
		}
		users = append(users, dtoUser)
	}
	return 0, users, nil
}
