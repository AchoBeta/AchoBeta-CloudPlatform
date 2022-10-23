package base

type Machine struct {
	Id        string `bson:"_id"`
	Ip        string `bson:"ip"`
	StartPort int    `bson:"startPort"`
	Memory    int64  `bson:"memory"`
	Core      int    `bson:"core"`
}
