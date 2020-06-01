package redis

import (
	"blog-go/config"
	"blog-go/model"
	"blog-go/tool"
	"fmt"
	pool "github.com/meitu/go-redis-pool"
	"gopkg.in/mgo.v2/bson"
)

var Client *pool.Pool

func init() {
	// oss.New() endpoint为空字符串时，将直接panic
	defer func() {
		if p := recover(); p != nil {
			fmt.Println("redis初始化失败: ", p)
		}
	}()
	poolInstance, e := pool.NewHA(config.RedisConfig)
	if e != nil {
		fmt.Println("redis连接失败: ", e)
		panic(e)
	}
	Client = poolInstance
	setupArticle()
}

// 初始化时将文章数据写入redis
func setupArticle() {
	var result []model.Article
	db := model.GetConn()
	defer db.Close()

	_ = db.C(model.CollectionArticle).Find(bson.M{
		"author":   "Calabash",
		"blogType": "public",
		"isActive": true,
	}).All(&result)

	for i := 0; i < len(result); i++ {
		var (
			nextBlog model.Article
			lastBlog model.Article
		)
		_ = db.C(model.CollectionArticle).Find(bson.M{"_id": bson.M{"$lt": result[i].ID}, "blogType": "public", "isActive": true}).Sort("-_id").Limit(1).One(&nextBlog)
		_ = db.C(model.CollectionArticle).Find(bson.M{"_id": bson.M{"$gt": result[i].ID}, "blogType": "public", "isActive": true}).Limit(1).One(&lastBlog)
		_ = Client.HMSet(tool.ObjectID2Str(result[i].ID), tool.StructToMap(model.ArticleWithMeta{
			Article:    result[i],
			NextBlogId: nextBlog.ID,
			LastBlogId: lastBlog.ID,
		}))
	}
}
