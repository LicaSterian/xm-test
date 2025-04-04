package models

type User struct {
	Username       string   `bson:"username"`
	HashedPassword string   `bson:"hashed_password"`
	Scopes         []string `bson:"scopes"`
}
