package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"tiktok/pkg/common"
	"tiktok/service"
)

var videoService service.VideoService

// PublishCJS
// 1. 获取客户端返回信息
// 2. 处理数据
// 3. 响应客户端
func PublishCJS(ctx *gin.Context) {
	//获取前端请求
	token := ctx.PostForm("token")
	if token == "" {
		ctx.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "没有获取到用户token",
		})
		log.Panicln("获取用户token失败")
		return
	}
	title := ctx.PostForm("title")
	if title == "" {
		log.Println("获取投稿作品标题失败")
		ctx.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "没有获取到投稿作品标题",
		})
		log.Panicln("获取投稿作品标题失败")
		return
	}
	data, err := ctx.FormFile("data")
	if err != nil {
		ctx.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "获取投稿视频失败",
		})
		log.Panicln("获取投稿视频失败", err)
		return
	}
	//处理数据
	filename, err := videoService.Publish(data, title, token)
	if err != nil {
		ctx.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "视频发布失败",
		})
		log.Panicln("视频发布失败")
		return
	}
	ctx.JSON(http.StatusOK, common.BaseResponse{
		StatusCode: 0,
		StatusMsg:  fmt.Sprint(filename, "视频发布成功"),
	})
}

func PublishListCJS(ctx *gin.Context) {
	hostToken := ctx.GetString("hostToken")
	if hostToken == "" {
		ctx.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "没有获取到用户token",
		})
		log.Panicln("获取用户token失败")
		return
	}
	var guestID uint
	if val, exists := ctx.Get("user_id"); !exists {
		ctx.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "没有获取到用户ID",
		})
		log.Panicln("userID获取失败")
		return
	} else {
		guestID = val.(uint)
	}
	//videoList, _, author := videoService.PublishList(guestID, hostToken)
	videos, err := videoService.PublishList(guestID, hostToken)
	if err != nil {
		ctx.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "获取发布作品详情时失败",
		})
		//log.Panicln("获取发布作品详情时失败", videoList, err)
		return
	}
	ctx.JSON(http.StatusOK, common.VideoListResponse{
		BaseResponse: common.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "请求发布列表成功!",
		},
		VideoList: videos,
	})
}

func FeedCJS(ctx *gin.Context) {
	//token := ctx.GetString("token")
	//isLogin := token != ""
	//ctx.GetString("latest_time")
}
