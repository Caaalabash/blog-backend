package model

import "gopkg.in/mgo.v2/bson"

const CollectionUser = "users"

type userInfo struct {
	twitter string
	github  string
	juejin  string
}

type User struct {
	ID       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	UserName string        `json:"userName" bson:"userName"`
	UserPwd  string        `json:"userPwd" bson:"userPwd"`
	Avatar   string        `json:"avatar" bson:"avatar"`
	UserInfo userInfo      `json:"userInfo" bson:"userInfo"`
}
