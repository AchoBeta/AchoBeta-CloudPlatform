package service

import (
	"cloud-platform/global"
	"cloud-platform/internal/base/cloud"
	"context"
	"os/exec"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetImages() (int8, []base.Image, error) {
	collection := global.GetMgoDb("abcp").Collection("image")
	cur, err := collection.Find(context.TODO(), bson.M{"isDelete": false})
	if err != nil {
		return 1, nil, err
	}
	images := []cloud.Image{}
	for cur.Next(context.TODO()) {
		image := cloud.Image{}
		err = cur.Decode(&image)
		if err != nil {
			return 2, nil, err
		}
		images = append(images, image)
	}
	return 0, images, nil
}

func GetImageInfo(imageId string, image *base.Image) (int8, error) {
	collection := global.GetMgoDb("abcp").Collection("image")
	res := collection.FindOne(context.TODO(), bson.M{"_id": imageId})
	if res.Err() != nil {
		if res.Err() == mongo.ErrNilDocument {
			return 1, res.Err()
		} else {
			return 2, res.Err()
		}
	}
	res.Decode(&image)
	return 0, nil
}

func DeleteImage(imageId string) (int8, error) {
	// TODO: 删除数据库
	collection := global.GetMgoDb("abcp").Collection("image")
	update := bson.M{"$set": bson.M{"isDelete": "true"}}
	_, err := collection.UpdateByID(context.TODO(), imageId, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 1, err
		} else {
			return 2, err
		}
	}
	return 0, nil
}

func PushDockerImage(imageName string) (int8, error) {
	_, err := exec.Command(base.DOCKER, base.IMAGE_PUSH, imageName).Output()
	if err != nil {
		return 1, err
	}
	return 0, err
}

func PushK8SImage(imageName string) (int8, error) {
	return 0, nil
}
