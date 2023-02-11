package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"net/http"
	"testing"
	"tiktok/mapper"
	"tiktok/pkg/middleware"
)

func TestVideoService_Publish(t *testing.T) {
	var data *multipart.FileHeader
	var token string
	var title string
	var videoService VideoService
	middleware.InitLog()
	mapper.InitOSS()
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
	engine.Run(":12345")
	//fmt.Println(data)

}
