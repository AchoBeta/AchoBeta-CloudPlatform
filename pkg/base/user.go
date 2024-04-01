package base

type User struct {
	Id         string   `bson:"_id"`
	Username   string   `bson:"username"`
	Password   string   `bson:"password"`
	Name       string   `bson:"name"`
	Pow        string   `bson:"pow"` // 权限：0 管理员，1 用户，2 未校验
	Containers []string `bson:"containers"`
}

type DTOUser struct {
	Id         string   `bson:"_id,omitempty"`
	Username   string   `bson:"username,omitempty"`
	Name       string   `bson:"name,omitempty"`
	Pow        string   `bson:"pow,omitempty"`
	Containers []string `bson:"containers,omitempty"`
}
