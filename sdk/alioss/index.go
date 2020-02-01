package alioss

import (
	"blog-go/config"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var Bucket *oss.Bucket

func init() {
	// oss.New() endpoint为空字符串时，将直接panic
	defer func() {
		if p := recover(); p != nil {
			fmt.Println("sdk/alioss 初始化失败，将导致上传接口不可用")
		}
	}()
	client, err := oss.New(config.AliOssConfig.EndPoint, config.AliOssConfig.Ak, config.AliOssConfig.Sk)
	if err != nil {
		fmt.Println("sdk/alioss 初始化client失败，将导致上传接口不可用")
		return
	}
	Bucket, err = client.Bucket("calabash-static")
	if err != nil {
		fmt.Println("sdk/alioss 初始化bucket失败，将导致上传接口不可用")
	}
}
