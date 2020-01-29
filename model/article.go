package model

import "gopkg.in/mgo.v2/bson"

const CollectionArticle = "articles"

type Article struct {
	ID          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	BlogTitle   string        `json:"blogTitle" bson:"blogTitle"`
	BlogDate    string        `json:"blogDate" bson:"blogDate"`
	BlogContent string        `json:"blogContent" bson:"blogContent"`
	BlogType    string        `json:"blogType" bson:"blogType"`
	Author      string        `json:"author" bson:"author"`
	IsActive    bool          `json:"isActive" bson:"isActive"`
}
