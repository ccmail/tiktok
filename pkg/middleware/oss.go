package middleware

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"os"
)

var (
	//once              sync.Once
	//ossServerInstance *OssServer
	ossConfig  *OssConfig
	Bucket     *oss.Bucket
	EndPoint   string
	BucketName string
)

type OssConfig struct {
	AccessKeyId     string `yaml:"AccessKeyId"`
	AccessKeySecret string `yaml:"AccessKeySecret"`
	BucketName      string `yaml:"BucketName"`
	EndPoint        string `yaml:"EndPoint"`
}

type OssServer struct {
}

// InitOSS 初始化OSS服务
func InitOSS() error {
	err := os.MkdirAll("videos/", 0750)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	yamlFile, err := os.ReadFile("./config/oss.yaml")
	//yamlFile, err := os.ReadFile("E:\\OneDrive\\MyCode\\Go\\TikTok\\config\\oss.yaml")
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

// InitOSSSupportTest 初始化OSS服务
func InitOSSSupportTest() error {
	err := os.MkdirAll("videos/", 0750)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	yamlFile, err := os.ReadFile("../config/oss.yaml")
	//yamlFile, err := os.ReadFile("E:\\OneDrive\\MyCode\\Go\\TikTok\\config\\oss.yaml")
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

// OssUploadFromPath  上传至云端Oss，返回url
func OssUploadFromPath(filename string, filepath string) (url string, err error) {
	err = Bucket.PutObjectFromFile("short_video/"+filename, filepath)
	if err != nil {
		return "", err
	}
	url = "https://" + BucketName + "." + EndPoint + "/short_video/" + filename
	return url, nil
}

func OssUploadFromReader(filename string, data io.Reader) (url string, err error) {
	err = Bucket.PutObject("short_video/"+filename, data)
	if err != nil {
		return "", err
	}
	url = "https://" + BucketName + "." + EndPoint + "/short_video/" + filename
	return url, nil
}
