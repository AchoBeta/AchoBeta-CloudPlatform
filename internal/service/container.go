package service

import (
	"cloud-platform/global"
	"cloud-platform/internal/base"
	"cloud-platform/internal/base/cloud"
	"cloud-platform/internal/base/constant"
	"context"
	"fmt"

	"os/exec"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// 创建 Docker 容器
func CreateDockerContainer(container *cloud.Container, user *base.User) (error, int8) {
	container.Ip = global.Machine.Ip
	cmd := splitContainerCmd(container)
	out, err := exec.Command(constant.DOCKER, cmd...).Output()
	if err != nil {
		return err, 1
	}
	container.Id = string(out[:len(out)-1])
	container.Username = "root"
	container.Password = createToken()
	_, err = exec.Command(constant.DOCKER, constant.CONTAINER_EXEC, "-d", container.Name,
		"/bin/sh", "/bin/passwd.sh", container.Password).Output()
	if err != nil {
		return err, 2
	}
	container.Status = 0
	container.StartTime = time.Now().Unix()
	collection := global.GetMgoDb("abcp").Collection("container")
	_, err = collection.InsertOne(context.TODO(), container)
	if err != nil {
		return err, 2
	}
	collection = global.GetMgoDb("abcp").Collection("user")
	filter := bson.M{"_id": user.Id}
	res := collection.FindOne(context.TODO(), filter)
	user1 := &base.User{}
	res.Decode(user1)
	update := bson.M{"$push": bson.M{"containers": container.Id}}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err, 3
	}
	return nil, 0
}

// 创建 K8S 容器
func CreateK8SContainer(container *cloud.Container) (error, int8) {
	return nil, 0
}

// 根据 id 获取容器信息
func GetContainer(containerId string, container *cloud.Container) (error, int8) {
	collection := global.GetMgoDb("abcp").Collection("container")
	filter := bson.M{"_id": containerId}
	res := collection.FindOne(context.TODO(), filter)
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return res.Err(), 1
		} else {
			return res.Err(), 2
		}
	}
	err := res.Decode(container)
	if err != nil {
		return err, 3
	}
	return nil, 0
}

func GetContainers(user *base.User) (error, int8, []cloud.Container) {
	collection := global.GetMgoDb("abcp").Collection("container")
	filter := bson.M{"_id": bson.M{"$in": user.Containers}}
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return err, 1, nil
	}
	defer cur.Close(context.TODO())
	containers := []cloud.Container{}
	for cur.Next(context.TODO()) {
		container := cloud.Container{}
		err = cur.Decode(&container)
		if err != nil {
			return err, 2, nil
		}
		containers = append(containers, container)
	}
	return nil, 0, containers
}

// 删除 Docker 容器
func RemoveDockerContainer(containerId string, userId string) (error, int8) {
	_, err := exec.Command(constant.DOCKER, constant.CONTAINER_RM, "-f", containerId).Output()
	if err != nil {
		return err, 1
	}
	collection := global.GetMgoDb("abcp").Collection("container")
	_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": containerId})
	if err != nil {
		return err, 2
	}
	collection = global.GetMgoDb("abcp").Collection("user")
	_, err = collection.UpdateByID(context.TODO(), userId, bson.M{"$pull": bson.M{"containers": containerId}})
	if err != nil {
		return err, 3
	}
	return nil, 0
}

// 删除 K8S 容器
func RemoveK8SContainer(containerId string) (error, int8) {
	return nil, 0
}

// 开启 Docker 容器
func StartDockerContainer(containerId string) (error, int8) {
	_, err := exec.Command(constant.DOCKER, constant.CONTAINER_START, containerId).Output()
	if err != nil {
		return err, 1
	}
	collection := global.GetMgoDb("abcp").Collection("container")
	update := bson.M{"$set": bson.M{"status": 0, "startTime": time.Now().Unix()}}
	_, err = collection.UpdateByID(context.TODO(), containerId, update)
	if err != nil {
		return err, 2
	}
	return nil, 0
}

// 开启 K8S 容器
func StartK8SContainer(containerId string) (error, int8) {
	return nil, 0
}

// 停止 Docker 容器
func StopDockerContainer(containerId string) (error, int8) {
	_, err := exec.Command(constant.DOCKER, constant.CONTAINER_STOP, containerId).Output()
	if err != nil {
		return err, 1
	}
	collection := global.GetMgoDb("abcp").Collection("container")
	update := bson.M{"$set": bson.M{"status": 1}}
	_, err = collection.UpdateByID(context.TODO(), containerId, update)
	if err != nil {
		return err, 2
	}
	return nil, 0
}

// 停止 K8S 容器
func StopK8SContainer(containerId string) (error, int8) {
	return nil, 0
}

// 重启 Docker 容器
func RestartDockerContainer(containerId string) (error, int8) {
	_, err := exec.Command(constant.DOCKER, constant.CONTAINER_RESTART, containerId).Output()
	if err != nil {
		return err, 1
	}
	collection := global.GetMgoDb("abcp").Collection("container")
	filter := bson.M{"containerId": containerId}
	update := bson.M{"$set": bson.M{"Status": 0}}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err, 2
	}
	return nil, 0
}

// 重启 K8S 容器
func RestartK8SContainer(containerId string) (error, int8) {
	return nil, 0
}

// 根据 Docker 容器制作镜像
func MakeDockerImage(containerId string, image *cloud.Image) (error, int8) {
	out, err := exec.Command(constant.DOCKER, constant.CONTAINER_COMMIT, "-a", image.Author,
		"-m", image.Desc, containerId, image.Name+":"+image.Tag).Output()
	if err != nil {
		return err, 1
	}
	image.Id = string(out[7:19])
	collection := global.GetMgoDb("abcp").Collection("image")
	_, err = collection.InsertOne(context.TODO(), image)
	if err != nil {
		return err, 2
	}
	return nil, 0
}

// 根据 K8S 容器制作镜像
func MakeK8SImage(containerId string, image *cloud.Image) (error, int8) {
	return nil, 0
}

// 获取 Docker 容器日志
func GetDockerContainerLog(containerId string) (error, int8, string) {
	out, err := exec.Command(constant.DOCKER, constant.CONTAINER_LOG, containerId).Output()
	if err != nil {
		return err, 1, ""
	}
	return nil, 0, string(out)
}

// 获取 K8S 容器日志
func GetK8SContainerLog(containerId string) (error, int8, string) {
	return nil, 0, ""
}

func splitContainerCmd(container *cloud.Container) []string {
	strs := []string{constant.CONTAINER_RUN, "-d", "--privileged=true", "--restart=always"}
	strs = append(strs, "--name", container.Name)
	if container.Param.Env != nil {
		for _, value := range container.Param.Env {
			strs = append(strs, "-e", value)
		}
	}
	container.Ports = global.Machine.StartPort
	strs = append(strs, "-p", fmt.Sprintf("%d:22", container.Ports))
	strs = append(strs, "-p", fmt.Sprintf("%d:23", container.Ports+1))
	for i := 0; i < 8; i++ {
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
