package model

import "gopkg.in/mgo.v2/bson"

const CollectionUser = "users"

type userInfo struct {
	Twitter string `json:"twitter" bson:"twitter"`
	Github  string `json:"github" bson:"github"`
	Juejin  string `json:"juejin" bson:"juejin"`
}

type User struct {
	ID       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	UserName string        `json:"userName" bson:"userName"`
	UserPwd  string        `json:"userPwd,omitempty" bson:"userPwd"`
	Avatar   string        `json:"avatar" bson:"avatar"`
	UserInfo userInfo      `json:"userInfo" bson:"userInfo"`
}
