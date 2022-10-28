package config

import (
	"blog-go/tool"
	pool "github.com/meitu/go-redis-pool"
	"os"
)

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Total   int         `json:"total"`
}

type OssConfig struct {
	Url        string
	EndPoint   string
	Ak         string
	Sk         string
	BucketName string
}

const FailedCode = 1
const SuccessCode = 0

var AppHost = tool.GetEnv("APP_HOST", "dockerhost")
var MongoURL = AppHost

var AliOssConfig = &OssConfig{
	Url:        "https://static.calabash.top/",
	EndPoint:   os.Getenv("ALI_ENDPOINT"),
	Ak:         os.Getenv("ALI_AK"),
	Sk:         os.Getenv("ALI_SK"),
	BucketName: os.Getenv("ALI_BUCKET"),
}

var RedisConfig = &pool.HAConfig{
	Master: AppHost + ":9000",
	Slaves: []string{
		AppHost + ":9001",
		AppHost + ":9002",
	},
}
