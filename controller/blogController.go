package controller

import (
	"blog-go/config"
	"blog-go/model"
	"github.com/kataras/iris/v12"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

type BlogController struct {
	collection string
}

type ArticleWithMeta struct {
	Article    model.Article `json:"article"`
	NextBlogId bson.ObjectId `json:"nextBlogId"`
	LastBlogId bson.ObjectId `json:"lastBlogId"`
}

func NewBlogController() BlogController {
	return BlogController{
		collection: model.CollectionArticle,
	}
}

// 获取文章列表，分为两种情况
// 1. 全量获取：所有文章所有内容
// 2. 列表获取：公开文章的标题以及日期
func (c *BlogController) GetArticles(ctx iris.Context) {
	var result []model.Article
	var query *mgo.Query
	pgS, _ := strconv.Atoi(ctx.URLParam("pgS"))
	pgN, _ := strconv.Atoi(ctx.URLParam("pgN"))
	skipCount := pgS * (pgN - 1)

	db := model.GetConn()
	defer db.Close()

	if ctx.URLParam("type") == "public" {
		query = db.C(c.collection).Find(bson.M{
			"author":   ctx.URLParam("author"),
			"blogType": "public",
		}).Select(bson.M{
			"blogDate":  1,
			"blogTitle": 1,
		})
	} else {
		query = db.C(c.collection).Find(bson.M{
			"author": ctx.URLParam("author"),
		})
	}
	e := query.Sort("-blogDate").Skip(skipCount).Limit(pgS).All(&result)
	if e != nil {
		_, _ = ctx.JSON(&config.Response{
			Code:    config.FailedCode,
			Message: e.Error(),
		})
	} else {
		_, _ = ctx.JSON(&config.Response{
			Code: config.SuccessCode,
			Data: result,
		})
	}
}

// 获取具体的某一篇文章，需要额外获取上一篇和下一篇的信息
// 下一篇 = <id
// 上一篇 = >id
func (c *BlogController) GetArticle(ctx iris.Context) {
	var currentBlog model.Article
	var nextBlog model.Article
	var lastBlog model.Article
	id := bson.ObjectIdHex(ctx.Params().Get("id"))

	db := model.GetConn()
	defer db.Close()

	e := db.C(c.collection).Find(bson.M{"_id": id}).One(&currentBlog)
	_ = db.C(c.collection).Find(bson.M{"_id": bson.M{"$lt": id}, "blogType": "public"}).Sort("-_id").Limit(1).One(&nextBlog)
	_ = db.C(c.collection).Find(bson.M{"_id": bson.M{"$gt": id}, "blogType": "public"}).Limit(1).One(&lastBlog)

	if e != nil {
		_, _ = ctx.JSON(&config.Response{
			Code:    config.FailedCode,
			Data:    nil,
			Message: e.Error(),
		})
	} else {
		_, _ = ctx.JSON(&config.Response{
			Code: config.SuccessCode,
			Data: ArticleWithMeta{
				currentBlog,
				nextBlog.ID,
				lastBlog.ID,
			},
		})
	}
}
