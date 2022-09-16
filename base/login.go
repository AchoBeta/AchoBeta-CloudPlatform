package base

type User struct {
	Id       int64  `bson:"_id"`
	Username string `bson:"username"`
	Password string `bson:"password"`
	/** todo */
}
