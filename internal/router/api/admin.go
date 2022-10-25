package api

import (
	"cloud-platform/global"
	"cloud-platform/internal/base"
	"cloud-platform/internal/base/config"
	"cloud-platform/internal/handle"
	"cloud-platform/internal/router"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	router.Register(func(router gin.IRoutes) {
		router.GET("/unverity_users", getUnVerityUsers)
		router.POST("/verity_user/:id", verityUser)
		router.POST("/set-admin/:id", setAdmin)
	}, router.V2)
}

func getUnVerityUsers(c *gin.Context) {
	r := handle.NewResponse(c)
	collection := global.GetMgoDb("abcp").Collection("user")
	res, err := collection.Find(context.TODO(), bson.M{"pow": config.TOURIST_POW})
	if err != nil {
		glog.Errorf("[db] get un-verity user error ! msg: %v\n", err.Error())
		r.Error(handle.INTERNAL_ERROR)
		return
	}
	var users []base.DTOUser
	for res.Next(context.TODO()) {
		var user base.DTOUser
		res.Decode(&user)
		users = append(users, user)
	}
	r.Success(users)
}

func verityUser(c *gin.Context) {
	r := handle.NewResponse(c)
	id := c.Query("id")
	collection := global.GetMgoDb("abcp").Collection("user")
	update := bson.M{"$set": bson.M{"pow": config.USER_POW}}
	_, err := collection.UpdateByID(context.TODO(), id, update)
	if err != nil {
		glog.Errorf("[db] get verity user error ! msg: %v\n", err.Error())
		r.Error(handle.INTERNAL_ERROR)
		return
	}
	r.Success(nil)
}

func setAdmin(c *gin.Context) {
	r := handle.NewResponse(c)
	id := c.Query("id")
	collection := global.GetMgoDb("abcp").Collection("user")
	update := bson.M{"$set": bson.M{"pow": 0}}
	_, err := collection.UpdateByID(context.TODO(), id, update)
	if err != nil {
		glog.Errorf("[db] get verity user error ! msg: %v\n", err.Error())
		r.Error(handle.INTERNAL_ERROR)
		return
	}
	r.Success(nil)
}
