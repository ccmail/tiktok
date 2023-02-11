package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"tiktok/pkg/common"
	"tiktok/service"
)

// PublishCJS
// 1. 获取客户端返回信息
// 2. 处理数据
// 3. 响应客户端
func PublishCJS(c *gin.Context) {
	//获取前端请求
	token := c.PostForm("token")
	if token == "" {
		log.Println("获取用户token失败")
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "没有获取到用户token",
		})
		return
	}
	title := c.PostForm("title")
	if title == "" {
		log.Println("获取投稿作品标题失败")
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "没有获取到投稿作品标题",
		})
		return
	}
	data, err := c.FormFile("data")
	if err != nil {
		log.Println("获取投稿视频失败")
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "获取投稿视频失败",
		})
		return
	}
	//处理数据
	var videoService service.VideoService
	filename, err := videoService.Publish(data, title, token)
	if err != nil {
		log.Println("视频发布失败", err)
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "视频发布失败",
		})
		return
	}
	c.JSON(http.StatusOK, common.BaseResponse{
		StatusCode: 0,
		StatusMsg:  fmt.Sprint(filename, "视频发布成功"),
	})
}
func FeedCJS() {

}
