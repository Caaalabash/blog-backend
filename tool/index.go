package tool

import (
	"encoding/json"
	"fmt"
	redisType "github.com/go-redis/redis/v7"
	pool "github.com/meitu/go-redis-pool"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/russross/blackfriday.v2"
	"html/template"
	"time"
	"unsafe"
)

func Str2Bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func ObjectID2Str(id bson.ObjectId) string {
	return fmt.Sprintf("%x", string(id))
}

func MarkdownToHtml(md string) template.HTML {
	b := Str2Bytes(md)
	b = blackfriday.Run(b, blackfriday.WithNoExtensions())
	return template.HTML(Bytes2Str(b))
}

func StructToMap(obj interface{}) map[string]interface{} {
	var m map[string]interface{}
	jsonBytes, _ := json.Marshal(obj)
	_ = json.Unmarshal(jsonBytes, &m)
	return m
}

// Redis限流
func IsActionAllowed(client *pool.Pool, userID string, actionKey string, period time.Duration, maxCount int) bool {
	key := fmt.Sprintf("blog-limit:%s:%s", userID, actionKey)
	window := period.Nanoseconds() / 1e6
	now := time.Now().UnixNano() / 1e6

	pipe, _ := client.Pipeline()
	pipe.ZAdd(key, &redisType.Z{
		Score:  float64(now),
		Member: now,
	})
	pipe.ZRemRangeByScore(key, "0", fmt.Sprintf("%v", now-window))
	pipe.ZCard(key)
	pipe.Expire(key, period+time.Second*1)

	res, err := pipe.Exec()
	if err != nil {
		return false
	}
	cmd, ok := res[2].(*redisType.IntCmd)
	if ok {
		return cmd.Val() <= int64(maxCount)
	}
	return false
}
