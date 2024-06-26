package exec

import (
	"cloud-platform/global"
	"cloud-platform/internal/base/cloud"
	"cloud-platform/internal/base/config"
	"cloud-platform/internal/base/constant"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/go-redis/redis"
	"github.com/golang/glog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v2"
)

func Init(path string) {
	// 日志启动要放在最开始
	readConfig(path)
	initMongo()
	initRedis()
	initMachineInfo()
	initBaseImage()
	global.HttpClient = http.DefaultClient
}

func readConfig(file string) {
	// 导入配置文件
	global.Config = &config.Server{}
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		panic(err.Error())
	}
	// 将配置文件读取到结构体中
	err = yaml.Unmarshal(yamlFile, global.Config)
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
		panic(err.Error())
	}
	// 检查连接
	err = global.Mgo.Ping(context.TODO(), nil)
	if err != nil {
		panic(err.Error())
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
		glog.Errorf("mongo list databases error: %s", err)
		return err
	}
	for _, db := range databases {
		if db == "abcp" {
			return err
		}
	}
	glog.Errorf("Database 'abcp' does not exist")
	// 或者创建数据库
	_, err = global.Mgo.Database("abcp").Collection("image").InsertOne(context.TODO(), bson.M{})
	if err != nil {
		glog.Errorf("mongo create database error: %s", err)
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
		panic(err.Error())
	}
}

// 初始化基础镜像
func initBaseImage() {
	collection := global.GetMgoDb("abcp").Collection("image")
	imageName := fmt.Sprintf("%s/abcp-base", global.Config.Docker.Hub.Host)
	filter := bson.M{"name": imageName}
	res := collection.FindOne(context.TODO(), filter)
	// 拉取远程镜像
	// glog.Infof("====== [cmd] pull base images ======")
	// _, err := exec.Command(base.DOCKER, base.IMAGE_PULL, imageName+":0.2").Output()
	// if err != nil {
	// 	glog.Errorf("[cmd] pull base images error ! msg: %s\n", err.Error())
	// }
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			out, err := exec.Command(constant.DOCKER, constant.IMAGES, imageName+":0.2").Output()
			if err != nil {
				glog.Errorf("[cmd] search base images error ! msg: %s\n", err.Error())
				return
			}
			fmt.Print(string(out))
			ss := strings.Split(string(out), "\n")
			ss = strings.Split(ss[1], " ")
			image := cloud.Image{
				Name:       imageName,
				Tag:        "0.2",
				Id:         ss[10],
				CreateTime: "",
				Size:       "",
				Author:     "abcp",
				Desc:       "base image; include ssh,scp; should bind port 22",
			}
			_, err = collection.InsertOne(context.TODO(), &image)
			if err != nil {
				glog.Errorf("[db] insert base images error ! msg: %s\n", err.Error())
				return
			}
		} else {
			glog.Error("[db] find base images error ! msg: %s\n", res.Err().Error())
		}
	}
}

// 初始化本地机器的信息
func initMachineInfo() {
	// 从数据库读取
	collection := global.GetMgoDb("abcp").Collection("machine")
	filter := bson.M{"_id": "1"}
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
				glog.Errorf("[db] insert machin info error ! msg: %s\n", err.Error())
			}
		}
		glog.Errorf("[db] find machine info error ! msg: %s\n", res.Err().Error())
		return
	}
	res.Decode(global.Machine)
}
