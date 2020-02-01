package config

import "os"

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
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
