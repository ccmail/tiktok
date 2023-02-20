package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"testing"
	"tiktok/config"
	"tiktok/pkg/middleware"
	"time"
)

func TestVideoService_Publish(t *testing.T) {
	var data *multipart.FileHeader
	var token string
	var title string
	var videoService VideoService
	config.InitLog()
	_ = middleware.InitOSS()
	engine := gin.Default()
	engine.POST("/douyin/publish/action/", func(context *gin.Context) {
		data, _ = context.FormFile("data")
		token = context.PostForm("token")
		title = context.PostForm("title")
		context.String(http.StatusOK, fmt.Sprint(" ", token, " ", title))
		publish, err := videoService.Publish(data, title, token)
		if err != nil {
			log.Println(err)
			t.Error("发布失败")
		} else {
			fmt.Println(publish)
		}
	})
	_ = engine.Run(":12345")
	//fmt.Println(data)

}
func TestVideoService_PublishList(t *testing.T) {
	_ = config.InitDBConnectorSupportTest()
	var v VideoService
	list, err := v.PublishList(1, "")
	if err != nil {
		t.Error("获取失败了")
	}
	for i := range list {
		fmt.Println(list[i])
		//fmt.Println(i)
	}
}

func TestVideoService_Feed(t *testing.T) {
	_ = config.InitDBConnectorSupportTest()
	var v VideoService
	feed, nextTime, err := v.Feed("", "")
	if err != nil {
		t.Error()
	}
	log.Println(feed, nextTime)
}
func TestSetName(t *testing.T) {
	//saveName := fmt.Sprint(parseToken.UserId, time.Now().UnixNano(), fileName)
	saveName := fmt.Sprint(1, time.Now().UnixNano(), "南科大xxx  sce.txt")
	savePath := filepath.Join("./videos/", saveName)

	println(savePath)

}
