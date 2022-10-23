package service

import (
	"CloudPlatform/global"
	"CloudPlatform/internal/base"
	"context"
	"os/exec"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetImages() (error, int8, []base.Image) {
	collection := global.GetMgoDb("abcp").Collection("image")
	cur, err := collection.Find(context.TODO(), bson.M{"isDelete": false})
	if err != nil {
		return err, 1, nil
	}
	images := []base.Image{}
	for cur.Next(context.TODO()) {
		image := base.Image{}
		err = cur.Decode(&image)
		if err != nil {
			return err, 2, nil
		}
		images = append(images, image)
	}
	return nil, 0, images
}

func GetImageInfo(imageId string, image *base.Image) (error, int8) {
	collection := global.GetMgoDb("abcp").Collection("image")
	res := collection.FindOne(context.TODO(), bson.M{"_id": imageId})
	if res.Err() != nil {
		if res.Err() == mongo.ErrNilDocument {
			return res.Err(), 1
		} else {
			return res.Err(), 2
		}
	}
	res.Decode(&image)
	return nil, 0
}

func DeleteImage(imageId string) (error, int8) {
	// TODO: 删除数据库
	collection := global.GetMgoDb("abcp").Collection("image")
	update := bson.D{{"$set", bson.D{{"isDelete", "true"}}}}
	_, err := collection.UpdateByID(context.TODO(), imageId, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return err, 1
		} else {
			return err, 2
		}
	}
	return err, 0
}

func PushDockerImage(imageName string) (error, int8) {
	_, err := exec.Command(base.DOCKER, base.IMAGE_PUSH, imageName).Output()
	if err != nil {
		return err, 1
	}
	return err, 0
}

func PushK8SImage(imageName string) (error, int8) {
	return nil, 0
}
