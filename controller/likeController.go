package controller

import (
	"log"
	"net/http"
	"strconv"
	"tiktok/pkg/common"
	"tiktok/service"

	"github.com/gin-gonic/gin"
)

// Like 点赞视频控制层
func Like(c *gin.Context) {
	// 由token获取user_id（jwt中间件完成）
	getUserId, _ := c.Get("user_id")
	var userId uint
	if v, ok := getUserId.(uint); ok {
		userId = v
	}

	actionTypeStr := c.Query("action_type")
	actionType, _ := strconv.ParseUint(actionTypeStr, 10, 10)
	videoIdStr := c.Query("video_id")
	videoId, _ := strconv.ParseUint(videoIdStr, 10, 10)

	// 函数调用及响应
	if actionType == 1 {
		err := service.LikeService(userId, uint(videoId), uint(actionType))
		if err != nil {
			log.Panicln("controller-Like: 点赞失败，", err)
			c.JSON(http.StatusBadRequest, common.BaseResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			return
		}
		log.Println("点赞成功")
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "点赞成功！",
		})
	} else if actionType == 2 {
		err := service.CancelLikeService(userId, uint(videoId))
		if err != nil {
			c.JSON(http.StatusBadRequest, common.BaseResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			return
		}
		log.Println("取消点赞成功")
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "取消点赞成功！",
		})
	} else {
		log.Printf("controller-Like: 操作失败：未定义的actionType%v。", actionType)
		return
	}

}

// LikeList 喜欢列表控制层
func LikeList(c *gin.Context) {
	getUserId, _ := c.Get("user_id") // token对应用户的id
	var userIdHost uint
	if v, ok := getUserId.(uint); ok {
		userIdHost = v
	}
	userIdStr := c.Query("user_id") // 目标用户的id
	userId, _ := strconv.ParseUint(userIdStr, 10, 10)
	userIdNew := uint(userId)
	if userIdNew == 0 {
		userIdNew = userIdHost
	}

	//函数调用及响应
	videoList, err := service.LikeListService(userIdNew)
	// log.Println("videoList:", videoList)
	returnList := service.FillInfo(videoList, userIdNew)
	// log.Println("returnList:", returnList)

	if err != nil {
		log.Panicln("controller-LikeList: 获取喜欢列表失败，", err)
		c.JSON(http.StatusBadRequest, common.LikeListResponse{
			BaseResponse: common.BaseResponse{
				StatusCode: 1,
				StatusMsg:  "获取喜欢列表失败",
			},
			VideoList: nil,
		})
		return
	} else {
		if len(returnList) == 0 {
			c.JSON(http.StatusOK, common.LikeListResponse{
				BaseResponse: common.BaseResponse{
					StatusCode: 0,
					StatusMsg:  "该用户未点赞任何视频",
				},
			})
		}
		c.JSON(http.StatusOK, common.LikeListResponse{
			BaseResponse: common.BaseResponse{
				StatusCode: 0,
				StatusMsg:  "获取喜欢列表成功",
			},
			VideoList: returnList,
		})
	}

}
