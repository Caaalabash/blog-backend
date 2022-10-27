package config

import (
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
const MongoURL = "dockerhost"

var AliOssConfig = &OssConfig{
	Url:        "https://static.calabash.top/",
	EndPoint:   os.Getenv("ALI_ENDPOINT"),
	Ak:         os.Getenv("ALI_AK"),
	Sk:         os.Getenv("ALI_SK"),
	BucketName: os.Getenv("ALI_BUCKET"),
}

var RedisConfig = &pool.HAConfig{
	Master: "dockerhost:9000",
	Slaves: []string{
		"dockerhost:9001",
		"dockerhost:9002",
	},
}
