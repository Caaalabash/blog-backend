package controller

import (
	"blog-go/config"
	"blog-go/model"
	"blog-go/redis"
	"blog-go/sdk/alioss"
	"blog-go/tool"
	"fmt"
	"github.com/kataras/iris/v12"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"regexp"
	"strconv"
)

type BlogController struct {
	collection string
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
			"isActive": true,
		}).Select(bson.M{
			"blogDate":  1,
			"blogTitle": 1,
		})
	} else {
		query = db.C(c.collection).Find(bson.M{
			"author":   ctx.URLParam("author"),
			"isActive": true,
		})
	}
	total, _ := query.Count()
	e := query.Sort("-blogDate").Skip(skipCount).Limit(pgS).All(&result)
	if e != nil {
		panic(e)
	} else {
		_, _ = ctx.JSON(&config.Response{
			Code:  config.SuccessCode,
			Data:  result,
			Total: total,
		})
	}
}

// 获取具体的某一篇文章，需要额外获取上一篇和下一篇的信息
// 下一篇 = <id
// 上一篇 = >id
func (c *BlogController) GetArticle(ctx iris.Context) {
	id := ctx.Params().Get("id")
	if cache := redis.Client.HGetAll(id).Val(); len(cache) != 0 {
		_, _ = ctx.JSON(&config.Response{
			Code: config.SuccessCode,
			Data: cache,
		})
		return
	}
	var currentBlog model.Article
	var nextBlog model.Article
	var lastBlog model.Article
	objectId := bson.ObjectIdHex(id)

	db := model.GetConn()
	defer db.Close()

	e := db.C(c.collection).Find(bson.M{"_id": objectId}).One(&currentBlog)
	_ = db.C(c.collection).Find(bson.M{"_id": bson.M{"$lt": objectId}, "blogType": "public", "isActive": true}).Sort("-_id").Limit(1).One(&nextBlog)
	_ = db.C(c.collection).Find(bson.M{"_id": bson.M{"$gt": objectId}, "blogType": "public", "isActive": true}).Limit(1).One(&lastBlog)
	_ = redis.Client.HMSet(id, tool.StructToMap(model.ArticleWithMeta{
		Article:    currentBlog,
		NextBlogId: nextBlog.ID,
		LastBlogId: lastBlog.ID,
	}))

	if e != nil {
		panic(e)
	} else {
		_, _ = ctx.JSON(&config.Response{
			Code: config.SuccessCode,
			Data: model.ArticleWithMeta{
				Article:    currentBlog,
				NextBlogId: nextBlog.ID,
				LastBlogId: lastBlog.ID,
			},
		})
	}
}

// 删除文章，逻辑删除
func (c *BlogController) DeleteArticle(ctx iris.Context) {
	id := bson.ObjectIdHex(ctx.Params().Get("id"))

	db := model.GetConn()
	defer db.Close()

	_, e := db.C(c.collection).Find(bson.M{"_id": id}).Apply(mgo.Change{
		Update: bson.M{"$set": bson.M{"isActive": false}},
	}, nil)
	_ = redis.Client.HSet(ctx.Params().Get("id"), "isActive", false)

	if e != nil {
		panic(e)
	} else {
		_, _ = ctx.JSON(&config.Response{
			Code:    config.SuccessCode,
			Message: "删除成功",
		})
	}
}

// 修改文章
func (c *BlogController) UpdateArticle(ctx iris.Context) {
	var body model.Article
	_ = ctx.ReadJSON(&body)
	id := bson.ObjectIdHex(ctx.Params().Get("id"))

	db := model.GetConn()
	defer db.Close()

	_, e := db.C(c.collection).Find(bson.M{"_id": id}).Apply(mgo.Change{
		Update: bson.M{
			"$set": bson.M{
				"blogType":    body.BlogType,
				"blogContent": body.BlogContent,
				"blogTitle":   body.BlogTitle,
			},
		},
	}, nil)
	_ = redis.Client.HMSet(ctx.Params().Get("id"), map[string]interface{}{
		"blogType":    body.BlogType,
		"blogContent": body.BlogContent,
		"blogTitle":   body.BlogTitle,
	})

	if e != nil {
		panic(e)
	} else {
		_, _ = ctx.JSON(&config.Response{
			Code:    config.SuccessCode,
			Message: "修改成功",
		})
	}
}

// 创建文章
func (c *BlogController) CreateArticle(ctx iris.Context) {
	body := model.Article{IsActive: true}
	_ = ctx.ReadJSON(&body)

	db := model.GetConn()
	defer db.Close()

	e := db.C(c.collection).Insert(body)

	if e != nil {
		panic(e)
	} else {
		_, _ = ctx.JSON(&config.Response{
			Code:    config.SuccessCode,
			Message: "发布成功",
		})
	}
}

// 图片上传
func (c *BlogController) UploadFile(ctx iris.Context) {
	file, info, err := ctx.FormFile("uploadfile")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	filename := "blog-media/file/file-" + info.Filename
	err = alioss.Bucket.PutObject(filename, file)
	if err != nil {
		panic(err)
	}
	_, _ = ctx.JSON(&config.Response{
		Code:    config.SuccessCode,
		Data:    config.AliOssConfig.Url + filename,
		Message: "上传成功",
	})
}

// 渲染列表页
func (c *BlogController) RenderList(ctx iris.Context) {
	var result []model.Article

	db := model.GetConn()
	defer db.Close()

	_ = db.C(c.collection).Find(bson.M{
		"author":   "Calabash",
		"blogType": "public",
		"isActive": true,
	}).Select(bson.M{"blogTitle": 1}).Sort("-blogDate").All(&result)

	ctx.ViewData("list", result)
	_ = ctx.View("index.html")
}

// 渲染文章
func (c *BlogController) RenderArticle(ctx iris.Context) {
	if match, _ := regexp.MatchString(`\d{14}`, ctx.Params().Get("id")); match {
		ctx.StatusCode(404)
		return
	}
	var result model.Article
	id := bson.ObjectIdHex(ctx.Params().Get("id"))

	db := model.GetConn()
	defer db.Close()

	_ = db.C(c.collection).Find(bson.M{"_id": id}).One(&result)

	ctx.ViewData("article", result)
	_ = ctx.View("article.html")
}

// 生成站点地图
func (c *BlogController) GetSiteMap(ctx iris.Context) {
	var result []model.Article

	db := model.GetConn()
	defer db.Close()

	_ = db.C(c.collection).Find(bson.M{"author": "Calabash"}).Select(bson.M{"_id": 1}).All(&result)

	text := "https://blog.calabash.top/Calabash\n"
	for _, v := range result {
		text += fmt.Sprintf("https://blog.calabash.top/Calabash/articles/%x", string(v.ID)) + "\n"
	}

	_, _ = ctx.Text(text)
}
