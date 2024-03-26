package router_test

import (
	"cloud-platform/global"
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestInsertOne(t *testing.T) {
	collection := global.GetMgoDb("test").Collection("test")
	stu := Student{
		"Jack", "123456",
	}
	result, err := collection.InsertOne(context.TODO(), stu)
	if err != nil {
		t.Errorf("result: %v\n", result)
		return
	}
	t.Logf("insert-id: %v\n", result.InsertedID)
}

func TestInsertMany(t *testing.T) {
	collection := global.GetMgoDb("test").Collection("test")
	stu1 := Student{
		"Jack", "11111",
	}
	stu2 := Student{
		"Mark", "2222",
	}
	result, err := collection.InsertMany(context.TODO(), []interface{}{stu1, stu2})
	if err != nil {
		t.Errorf("result: %v\n", result)
		return
	}
	t.Logf("result-issue: %v\n", result.InsertedIDs)
}

func TestFindOne(t *testing.T) {
	collection := global.GetMgoDb("test").Collection("test")
	filter := bson.D{{"name", "Jack"}}
	var result Student
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		t.Errorf("result: %v\n", err)
		return
	}
	t.Logf("result: %v", result)
}

func TestFindMany(t *testing.T) {
	collection := global.GetMgoDb("test").Collection("test")
	filter := bson.D{{"name", "Jack"}}
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		t.Errorf("result: %v", err)
		return
	}
	for cur.Next(context.TODO()) {
		var stu Student
		cur.Decode(&stu)
		t.Logf("result: %v\n", stu)
	}
	if cur.Err() != nil {
		t.Error(cur.Err())
	}
	cur.Close(context.TODO())
}

func TestUpdateOne(t *testing.T) {
	collection := global.GetMgoDb("test").Collection("test")
	filter := bson.D{{"name", "Mark"}}
	update := bson.D{{"$set", bson.D{{"id", "3333"}}}}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("result: %v-%v\n", res.MatchedCount, res.ModifiedCount)
}

func TestUpdateMany(t *testing.T) {
	collection := global.GetMgoDb("test").Collection("test")
	filter := bson.D{{"name", "Jack"}}
	update := bson.D{{"$set", bson.D{{"name", "Tom"}}}}
	res, err := collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("result: %v-%v\n", res.MatchedCount, res.UpsertedCount)
}

func TestDeleteOne(t *testing.T) {
	collection := global.GetMgoDb("test").Collection("test")
	filter := bson.D{{"name", "Mark"}}
	res, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("result: %v\n", res.DeletedCount)
}

func TestDeleteMany(t *testing.T) {
	collection := global.GetMgoDb("test").Collection("test")
	filter := bson.D{{"name", "Tom"}}
	res, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("result: %v\n", res.DeletedCount)
}

func TestDeleteCollection(t *testing.T) {
	collection := global.GetMgoDb("test").Collection("test")
	err := collection.Drop(context.TODO())
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("success")
}
