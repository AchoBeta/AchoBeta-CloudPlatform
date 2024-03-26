package cloud

type Image struct {
	Id         string `bson:"_id"`
	Name       string `bson:"name"`
	Tag        string `bson:"tag"`
	Author     string `bson:"author"`
	Size       string `bson:"size"`
	Desc       string `bson:"describe"`
	CreateTime string `bson:"createTime"`
	IsDelete   bool   `bson:"isDelete"`
}
