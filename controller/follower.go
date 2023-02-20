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

func Follow(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == "" {
		ctx.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "没有获取到用户token",
		})
		log.Panicln("获取用户token失败")
		return
	}
	var guestID uint
	guestIDStr := ctx.Query("to_user_id")
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
	//true表示关注, false表示取关
	var isConcern bool
	if op := ctx.Query("action_type"); op == "" {
		ctx.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "关注/取关操作失败, 请重试",
		})
		log.Panicln("获取关注操作状态时失败")
		return
	} else {
		isConcern = op[0] == '1'
	}
	err := service.Follow(token, guestID, isConcern)
	if err != nil {
		ctx.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  fmt.Sprint("关注/取关失败", err),
		})
		log.Panicln("关注/取关失败, 注意检查service.Follow的逻辑")
		return
	}
	ctx.JSON(http.StatusOK, common.BaseResponse{
		StatusCode: 0,
		StatusMsg:  "操作成功!",
	})
}

// FollowList 关注列表
//与发布列表类似, 用户未登录应该也可以查看某位作者关注的人, 这里先做好准备, 目前走jwt鉴权, 为以后更新准备, 先做好token为空的判断
func FollowList(ctx *gin.Context) {
	hostToken := ctx.Query("token")
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
	userInfoList, err := service.FollowList(hostToken, guestID)
	if err != nil {
		ctx.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "获取关注列表失败!",
		})
		log.Panicln("获取关注列表时失败")
		return
	}
	ctx.JSON(http.StatusOK, common.UserInfoListBaseResp{
		BaseResponse: common.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "请求发布列表成功!",
		},
		UserList: userInfoList,
	})
}

// FollowerList 粉丝列表
//逻辑类似关注列表
func FollowerList(ctx *gin.Context) {
	hostToken := ctx.Query("token")
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
	userInfoList, err := service.FollowerList(hostToken, guestID)
	if err != nil {
		ctx.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "获取粉丝列表失败!",
		})
		log.Panicln("获取粉丝列表时失败")
		return
	}
	ctx.JSON(http.StatusOK, common.UserInfoListBaseResp{
		BaseResponse: common.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "请求粉丝列表成功!",
		},
		UserList: userInfoList,
	})
}

// FriendList
//好友列表, 查询逻辑是互相关注
func FriendList(ctx *gin.Context) {
	hostToken := ctx.Query("token")
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
	userInfoList, err := service.FriendList(hostToken, guestID)
	if err != nil {
		ctx.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "获取好友列表失败!",
		})
		log.Panicln("获取好友列表时失败")
		return
	}
	ctx.JSON(http.StatusOK, common.UserInfoListBaseResp{
		BaseResponse: common.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "请求好友列表成功!",
		},
		UserList: userInfoList,
	})
}
