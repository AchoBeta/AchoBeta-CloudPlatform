package service

import (
	"cloud-platform/global"
	"cloud-platform/internal/base"
	"cloud-platform/internal/base/cloud"
	"cloud-platform/internal/base/constant"
	commonx "cloud-platform/internal/pkg/common"
	"context"
	"fmt"
	"os/exec"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/text/encoding/simplifiedchinese"
)

// 创建容器
func CreateContainer(token string, container *cloud.Container, user *base.User) (int8, error) {
	container.Ip = global.Machine.Ip
	container.Username = "root"
	container.Password = createToken()
	if global.Config.App.Engine == "docker" {
		cmd := splitContainerCmd(container)
		out, err := exec.Command(constant.DOCKER, cmd...).Output()
		if err != nil {
			return 1, err
		}
		container.Id = string(out[:len(out)-1])
		_, err = exec.Command(constant.DOCKER, constant.CONTAINER_EXEC, "-d", container.Name,
			"/bin/sh", "/bin/passwd.sh", container.Password).Output()
		if err != nil {
			return 1, err
		}
	} else {
		cmd := splitContainerCmd(container)
		out, err := exec.Command(constant.K8S, cmd...).Output()
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
	cmd1 := global.Rdb.Set(fmt.Sprintf(base.TOKEN, token), str, 30*time.Minute)
	if cmd1.Err() != nil {
		return 4, cmd1.Err()
	}
	return 0, nil
}

// 根据 id 获取容器信息
func GetContainer(containerId string, container *cloud.Container) (int8, error) {
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

func GetContainers(user *base.User) (int8, []cloud.Container, error) {
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
		_, err = exec.Command(constant.DOCKER, constant.CONTAINER_RM, "-f", containerId).Output()
	} else {
		_, err = exec.Command(constant.K8S, constant.K8S_DELETE, containerId).Output()
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
	res := global.Rdb.Set(fmt.Sprintf(base.TOKEN, token), str, 30*time.Minute)
	if res.Err() != nil {
		return 4, err
	}
	return 0, nil
}

// 开启 Docker 容器
func StartDockerContainer(containerId string) (int8, error) {
	_, err := exec.Command(constant.DOCKER, constant.CONTAINER_START, containerId).Output()
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
	_, err := exec.Command(constant.DOCKER, constant.CONTAINER_STOP, containerId).Output()
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
	fmt.Println(containerId)
	_, err := exec.Command(constant.DOCKER, constant.CONTAINER_RESTART, containerId).Output()
	if err != nil {
		return 1, err
	}
	collection := global.GetMgoDb("abcp").Collection("container")
	// filter := bson.M{"containerId": containerId}
	update := bson.M{"$set": bson.M{"status": 0}}
	_, err = collection.UpdateByID(context.TODO(), containerId, update)
	if err != nil {
		return 2, err
	}
	return 0, nil
}

// 根据 Docker 容器制作镜像
func MakeDockerImage(containerId string, image *cloud.Image) (int8, error) {
	collection := global.GetMgoDb("abcp").Collection("image")
	filter := bson.M{"name": image.Name, "tag": image.Tag}
	res := collection.FindOne(context.TODO(), filter)
	if res.Err() != nil && res.Err() != mongo.ErrNoDocuments {
		return 2, res.Err()
	} else if res.Err() == nil {
		image1 := cloud.Image{}
		res.Decode(&image1)
		if !image1.IsDelete {
			return 3, nil
		}
	}
	out, err := exec.Command(constant.DOCKER, constant.CONTAINER_COMMIT, "-a", image.Author,
		"-m", image.Desc, containerId, image.Name+":"+image.Tag).Output()
	if err != nil {
		return 1, err
	}
	image.Id = string(out[7:19])

	_, err = collection.InsertOne(context.TODO(), image)
	if err != nil {
		return 2, err
	}
	return 0, nil
}

// 根据 K8S 容器制作镜像
func MakeK8SImage(containerId string, image *cloud.Image) (int8, error) {
	return 0, nil
}

// 获取 Docker 容器日志
func GetContainerLog(containerId string) (int8, string, error) {
	var out []byte
	var err error
	fmt.Print(containerId)
	if global.Config.App.Engine == "docker" {
		// print(exec.Command(constant.DOCKER, fmt.Sprintf("%s %s", constant.CONTAINER_LOG, containerId)))
		out, err = exec.Command(constant.DOCKER, constant.CONTAINER_LOG, containerId).Output()
	} else {
		out, err = exec.Command(constant.DOCKER, constant.K8S_LOG, containerId).Output()
	}
	if err != nil {
		return 1, "", err
	}
	fmt.Print(simplifiedchinese.GBK.NewDecoder().Bytes(out))
	return 0, string(out), nil
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
