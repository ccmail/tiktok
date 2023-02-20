package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"tiktok/pkg/common"
	"tiktok/service"
)

var videoService service.VideoService

// Publish
// 1. 获取客户端返回信息
// 2. 处理数据
// 3. 响应客户端
func Publish(ctx *gin.Context) {
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

func PublishList(ctx *gin.Context) {
	hostToken := ctx.Query("token")
	//假设用户处于未登录状态, 也应该可以访问某位作者的发布列表
	if hostToken == "" {
		log.Println("没有获取到用户token")
		log.Println(hostToken)
	}
	var guestID uint
	guestIDStr := ctx.Query("user_id")
	if atoi, err := strconv.Atoi(guestIDStr); err != nil || guestIDStr == "" {
		ctx.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "没有获取到用户ID",
		})
		log.Panicln("userID获取失败")
		return
	} else {
		guestID = uint(atoi)
	}

	videos, err := videoService.PublishList(guestID, hostToken)
	if err != nil {
		ctx.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "获取发布作品详情时失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, common.VideoListBaseResp{
		BaseResponse: common.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "请求发布列表成功!",
		},
		VideoList: videos,
	})
}

func Feed(ctx *gin.Context) {
	//token := ctx.GetString("token")
	token := ctx.Query("token")
	fmt.Println(token)
	//strLastTime := ctx.GetString("latest_time")
	strLastTime := ctx.Query("latest_time")
	videoResp, nextTime, err := videoService.Feed(token, strLastTime)
	if err != nil || len(videoResp) == 0 || nextTime <= 0 {
		ctx.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  fmt.Sprint(err),
		})
		return
	}

	ctx.JSON(http.StatusOK, common.FeedVideoListBaseResp{
		BaseResponse: common.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "获取Feed流成功",
		},
		NextTime:  nextTime,
		VideoList: videoResp,
	})

}
