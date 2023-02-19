package service

import (
	"CloudPlatform/global"
	"CloudPlatform/internal/base"
	commonx "CloudPlatform/pkg/common"
	"cloud-platform/internal/base/cloud"
	"context"
	"fmt"
	"os/exec"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// 创建容器
func CreateContainer(token string, container *base.Container, user *base.User) (int8, error) {
	container.Ip = global.Machine.Ip
	container.Username = "root"
	container.Password = createToken()
	if global.Config.App.Engine == "docker" {
		cmd := splitContainerCmd(container)
		out, err := exec.Command(base.DOCKER, cmd...).Output()
		if err != nil {
			return 1, err
		}
		container.Id = string(out[:len(out)-1])
		_, err = exec.Command(base.DOCKER, base.CONTAINER_EXEC, "-d", container.Name,
			"/bin/sh", "/bin/passwd.sh", container.Password).Output()
		if err != nil {
			return 1, err
		}
	} else {
		cmd := splitContainerCmd(container)
		out, err := exec.Command(base.K8S, cmd...).Output()
		if err != nil {
			return 1, err
		}
		container.Id = string(out[:len(out)-1])
	}
	container.Status = 0
	container.StartTime = time.Now().Unix()
	// 添加到数据库
	collection := global.GetMgoDb("abcp").Collection("container")
	_, err := collection.InsertOne(context.TODO(), container)
	if err != nil {
		return 2, err
	}
	// user 表更新并同步到缓存
	collection = global.GetMgoDb("abcp").Collection("user")
	update := bson.M{"$push": bson.M{"containers": container.Id}}
	_, err = collection.UpdateByID(context.TODO(), user.Id, update)
	if err != nil {
		return 3, err
	}
	user.Containers = append(user.Containers, container.Id)
	str, _ := commonx.StructToJson(&user)
	cmd1 := global.Rdb.Set(context.TODO(), fmt.Sprintf(base.TOKEN, token), str, 30*time.Minute)
	if cmd1.Err() != nil {
		return 4, cmd1.Err()
	}
	return 0, nil
}

// 根据 id 获取容器信息
func GetContainer(containerId string, container *base.Container) (int8, error) {
	collection := global.GetMgoDb("abcp").Collection("container")
	filter := bson.M{"_id": containerId}
	res := collection.FindOne(context.TODO(), filter)
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return 1, res.Err()
		} else {
			return 2, res.Err()
		}
	}
	err := res.Decode(container)
	if err != nil {
		return 3, err
	}
	return 0, nil
}

func GetContainers(user *base.User) (int8, []base.Container, error) {
	collection := global.GetMgoDb("abcp").Collection("container")
	filter := bson.M{"_id": bson.M{"$in": user.Containers}}
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return 1, nil, err
	}
	defer cur.Close(context.TODO())
	containers := []cloud.Container{}
	for cur.Next(context.TODO()) {
		container := cloud.Container{}
		err = cur.Decode(&container)
		if err != nil {
			return 2, nil, err
		}
		containers = append(containers, container)
	}
	return 0, containers, nil
}

// 删除容器
func RemoveContainer(token string, containerId string, user *base.User) (int8, error) {
	var err error
	if global.Config.App.Engine == "docker" {
		_, err = exec.Command(base.DOCKER, base.CONTAINER_RM, "-f", containerId).Output()
	} else {
		_, err = exec.Command(base.K8S, base.K8S_DELETE, containerId).Output()
	}
	if err != nil {
		return 1, err
	}
	collection := global.GetMgoDb("abcp").Collection("container")
	_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": containerId})
	if err != nil {
		return 2, err
	}
	collection = global.GetMgoDb("abcp").Collection("user")
	_, err = collection.UpdateByID(context.TODO(), user.Id, bson.M{"$pull": bson.M{"containers": containerId}})
	if err != nil {
		return 3, err
	}
	// 更新缓存
	index := len(user.Containers)
	for i, v := range user.Containers {
		if v == containerId {
			index = i
			break
		}
	}
	user.Containers = append(user.Containers[:index], user.Containers[index+1:]...)
	str, _ := commonx.StructToJson(&user)
	res := global.Rdb.Set(context.TODO(), fmt.Sprintf(base.TOKEN, token), str, 30*time.Minute)
	if res.Err() != nil {
		return 4, err
	}
	return 0, nil
}

// 开启 Docker 容器
func StartDockerContainer(containerId string) (int8, error) {
	_, err := exec.Command(base.DOCKER, base.CONTAINER_START, containerId).Output()
	if err != nil {
		return 1, err
	}
	collection := global.GetMgoDb("abcp").Collection("container")
	update := bson.M{"$set": bson.M{"status": 0, "startTime": time.Now().Unix()}}
	_, err = collection.UpdateByID(context.TODO(), containerId, update)
	if err != nil {
		return 2, err
	}
	return 0, nil
}

// 停止 Docker 容器
func StopDockerContainer(containerId string) (int8, error) {
	_, err := exec.Command(base.DOCKER, base.CONTAINER_STOP, containerId).Output()
	if err != nil {
		return 1, err
	}
	collection := global.GetMgoDb("abcp").Collection("container")
	update := bson.M{"$set": bson.M{"status": 1}}
	_, err = collection.UpdateByID(context.TODO(), containerId, update)
	if err != nil {
		return 2, err
	}
	return 0, nil
}

// 重启 Docker 容器
func RestartDockerContainer(containerId string) (int8, error) {
	_, err := exec.Command(base.DOCKER, base.CONTAINER_RESTART, containerId).Output()
	if err != nil {
		return 1, err
	}
	collection := global.GetMgoDb("abcp").Collection("container")
	filter := bson.M{"containerId": containerId}
	update := bson.M{"$set": bson.M{"Status": 0}}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return 2, err
	}
	return 0, nil
}

// 根据 Docker 容器制作镜像
func MakeDockerImage(containerId string, image *base.Image) (int8, error) {
	out, err := exec.Command(base.DOCKER, base.CONTAINER_COMMIT, "-a", image.Author,
		"-m", image.Desc, containerId, image.Name+":"+image.Tag).Output()
	if err != nil {
		return 1, err
	}
	image.Id = string(out[7:19])
	collection := global.GetMgoDb("abcp").Collection("image")
	_, err = collection.InsertOne(context.TODO(), image)
	if err != nil {
		return 2, err
	}
	return 3, nil
}

// 根据 K8S 容器制作镜像
func MakeK8SImage(containerId string, image *base.Image) (int8, error) {
	return 0, nil
}

// 获取 Docker 容器日志
func GetContainerLog(containerId string) (int8, string, error) {
	var out []byte
	var err error
	if global.Config.App.Engine == "docker" {
		out, err = exec.Command(base.DOCKER, base.CONTAINER_LOG, containerId).Output()
	} else {
		out, err = exec.Command(base.DOCKER, base.K8S_LOG, containerId).Output()
	}
	if err != nil {
		return 1, "", err
	}
	return 0, string(out), nil
}

func splitContainerCmd(container *base.Container) []string {
	strs := []string{base.CONTAINER_RUN, "-d", "--privileged=true", "--restart=always"}
	strs = append(strs, "--name", container.Name)
	if container.Param.Env != nil {
		for _, value := range container.Param.Env {
			strs = append(strs, "-e", value)
		}
	}
	container.Ports = global.Machine.StartPort
	container.Param.Ports = append([]int{22, 23}, container.Param.Ports...)
	for i := 0; i < 10; i++ {
		if container.Param.Ports != nil && i < len(container.Param.Ports) {
			strs = append(strs, "-p", fmt.Sprintf("%d:%d", container.Ports+i+2, container.Param.Ports[i]))
		} else {
			strs = append(strs, "-p", fmt.Sprintf("%d:%d", container.Ports+i+2, container.Ports+i+1))
		}
	}
	global.Machine.StartPort += 10
	if container.Param.HostName != "" {
		strs = append(strs, "-h", container.Param.HostName)
	}
	if container.Image != "" {
		strs = append(strs, container.Image)
	}
	return strs
}
