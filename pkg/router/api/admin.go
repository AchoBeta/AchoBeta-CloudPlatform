package api

import (
	"cloud-platform/global"
	"cloud-platform/pkg/base"
	"cloud-platform/pkg/base/config"
	"cloud-platform/pkg/handle"
	"cloud-platform/pkg/router/manager"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/route"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	manager.RouteHandler.RegisterRouter(manager.LEVEL_V2, func(router *route.RouterGroup) {
		router.GET("/unverity_users", getUnVerityUsers)
		router.POST("/verity_user/:id", verityUser)
		router.GET("/set-admin/:id", setAdmin)
	})
}

func getUnVerityUsers(ctx context.Context, c *app.RequestContext) {
	r := handle.NewResponse(c)
	collection := global.GetMgoDb("abcp").Collection("user")
	res, err := collection.Find(context.TODO(), bson.M{"pow": config.TOURIST_POW})
	if err != nil {
		global.Logger.Errorf("[db] get un-verity user error ! msg: %v\n", err.Error())
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

func verityUser(ctx context.Context, c *app.RequestContext) {
	r := handle.NewResponse(c)
	id := c.Query("id")
	collection := global.GetMgoDb("abcp").Collection("user")
	update := bson.M{"$set": bson.M{"pow": config.USER_POW}}
	_, err := collection.UpdateByID(context.TODO(), id, update)
	if err != nil {
		global.Logger.Errorf("[db] get verity user error ! msg: %v\n", err.Error())
		r.Error(handle.INTERNAL_ERROR)
		return
	}
	r.Success(nil)
}

func setAdmin(ctx context.Context, c *app.RequestContext) {
	r := handle.NewResponse(c)
	id := c.Query("id")
	collection := global.GetMgoDb("abcp").Collection("user")
	update := bson.M{"$set": bson.M{"pow": 0}}
	_, err := collection.UpdateByID(context.TODO(), id, update)
	if err != nil {
		global.Logger.Errorf("[db] get verity user error ! msg: %v\n", err.Error())
		r.Error(handle.INTERNAL_ERROR)
		return
	}
	r.Success(nil)
}
