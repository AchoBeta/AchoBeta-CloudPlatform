package load

import (
	"cloud-platform/global"
	"cloud-platform/pkg/base/cloud"
	"cloud-platform/pkg/base/config"
	"cloud-platform/pkg/base/constant"
	"flag"

	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const FILE_PATH = "./config.yaml"

func Init() {
	// 日志启动要放在最开始
	readConfig()
	initLog()
	initMongo()
	initRedis()
	initMachineInfo()
	initBaseImage()
}

func initLog() {
	logFilePath := flag.String("l", global.Config.Options.LogFilePath, "log file path")
	flag.Parse()
	InitLog(*logFilePath)
}

func readConfig() {
	//导入配置文件
	global.Config = &config.Server{}
	viper.SetConfigFile(FILE_PATH)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err.Error())
	}
	//将配置文件读取到结构体中
	err = viper.Unmarshal(global.Config)
	if err != nil {
		panic(err.Error())
	}
}

func initMongo() {
	var err error
	credential := options.Credential{
		Username:   global.Config.Db.Mongo.Username,
		Password:   global.Config.Db.Mongo.Password,
		AuthSource: global.Config.Db.Mongo.AuthSource,
	}
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d",
		global.Config.Db.Mongo.Address, global.Config.Db.Mongo.Port)).SetAuth(credential)
	// 连接到MongoDB
	global.Mgo, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		hlog.Errorf("mongo connect error: %s", err)
		return
	}
	// 检查连接
	err = global.Mgo.Ping(context.TODO(), nil)
	if err != nil {
		hlog.Errorf("mongo ping error: %s", err)
		return
	}
	// 检查所需要的数据库是否存在

	err = checkMongoDb()
	if err != nil {
		return
	}
}

func checkMongoDb() error {
	// 检查所需要的数据库是否存在
	databases, err := global.Mgo.ListDatabaseNames(context.TODO(), bson.M{})
	if err != nil {
		hlog.Errorf("mongo list databases error: %s", err)
		return err
	}
	for _, db := range databases {
		if db == "abcp" {
			return err
		}
	}
	hlog.Errorf("Database 'abcp' does not exist")
	// 或者创建数据库
	_, err = global.Mgo.Database("abcp").Collection("image").InsertOne(context.TODO(), bson.M{})
	if err != nil {
		hlog.Errorf("mongo create database error: %s", err)
		return err
	}
	return nil
}

func initRedis() {
	global.Rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.Config.Db.Redis.Address,
			global.Config.Db.Redis.Port),
		Password: global.Config.Db.Redis.Password,
		DB:       global.Config.Db.Redis.Db,
	})

	_, err := global.Rdb.Ping().Result()
	if err != nil {
		hlog.Errorf("redis connect fail! message: %s\n", err.Error())
	}
}

// 初始化基础镜像
func initBaseImage() {
	collection := global.GetMgoDb("abcp").Collection("image")
	imageName := fmt.Sprintf("%s/abcp_base", global.Config.Docker.Hub.Host)
	filter := bson.D{{Key: "name", Value: imageName}}
	res := collection.FindOne(context.TODO(), filter)
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			// 拉取远程镜像
			hlog.Infof("====== [cmd] pull base images ======")
			_, err := exec.Command(constant.DOCKER, constant.IMAGE_PULL, imageName+":0.1").Output()
			if err != nil {
				hlog.Errorf("[cmd] pull base images error ! msg: %s\n", err.Error())
			}
			out, err := exec.Command(constant.DOCKER, constant.IMAGES, imageName+"0.1").Output()
			if err != nil {
				hlog.Errorf("[cmd] search base images error ! msg: %s\n", err.Error())
				return
			}
			r := regexp.MustCompile(`[^\\s]+`)
			ss := r.FindAllString(strings.Split(string(out), "\n")[1], -1)
			fmt.Print(ss)
			image := cloud.Image{
				Name:       ss[0],
				Tag:        ss[1],
				Id:         ss[2],
				CreateTime: ss[3],
				Size:       ss[4],
				Author:     "abcp",
				Desc:       "base image; include ssh,scp; should bind port 22",
			}
			_, err = collection.InsertOne(context.TODO(), &image)
			if err != nil {
				hlog.Errorf("[db] insert base images error ! msg: %s\n", err.Error())
				return
			}
		} else {
			hlog.Error("[db] find base images error ! msg: %s\n", res.Err().Error())
		}
	}
}

// 初始化本地机器的信息
func initMachineInfo() {
	// 从数据库读取
	collection := global.GetMgoDb("abcp").Collection("machine")
	filter := bson.D{{Key: "_id", Value: "1"}}
	res := collection.FindOne(context.TODO(), filter)
	global.Machine = &cloud.Machine{}
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			// 初始化本机信息
			global.Machine.Id = "1"
			global.Machine.Ip = "127.0.0.1"  // TODO
			global.Machine.StartPort = 10000 // TODO
			global.Machine.Memory = 100      // TODO: 内存
			global.Machine.Core = 8          // TODO: 核心数
			_, err := collection.InsertOne(context.TODO(), global.Machine)
			if err != nil {
				hlog.Errorf("[db] insert machin info error ! msg: %s\n", err.Error())
			}
		}
		hlog.Errorf("[db] find machine info error ! msg: %s\n", res.Err().Error())
		return
	}
	res.Decode(global.Machine)
}
