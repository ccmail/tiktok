package mapper

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"gopkg.in/yaml.v2"
)

var (
	once                 sync.Once
	ossServerInstance    *OssServer
	ossConfig            *OssConfig
	Bucket               *oss.Bucket
	BucketName, EndPoint string
)

type OssConfig struct {
	AccessKeyId     string `yaml:"AccessKeyId"`
	AccessKeySecret string `yaml:"AccessKeySecret"`
	BucketName      string `yaml:"BucketName"`
	EndPoint        string `yaml:"EndPoint"`
}

type OssServer struct {
}

func NewOssServer() *OssServer {
	once.Do(func() {
		ossServerInstance = &OssServer{}
	})
	return ossServerInstance
}

// 初始化OSS服务
func InitOSS() error {
	err := os.MkdirAll("videos/", 0750)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	// yamlFile, err := os.ReadFile("./config/oss.yaml")
	yamlFile, err := os.ReadFile("D:\\Documents\\Golang Projects\\tiktok_ours\\config\\oss.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &ossConfig)
	BucketName, EndPoint = ossConfig.BucketName, ossConfig.EndPoint
	if err != nil {
		fmt.Println(err.Error())
	}

	client, err := oss.New(ossConfig.EndPoint, ossConfig.AccessKeyId, ossConfig.AccessKeySecret)
	if err != nil {
		return err
	}
	if ossConfig.BucketName != "" {
		Bucket, err = client.Bucket(ossConfig.BucketName)
		if err != nil {
			return err
		}
	}
	return nil
}
