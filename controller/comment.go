package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"tiktok/pkg/common"
	"tiktok/pkg/util"
	"tiktok/service"

	"github.com/gin-gonic/gin"
)

// Comment  评论操作控制层
func Comment(c *gin.Context) {
	//1 数据处理,这里
	userIdTemp, _ := c.Get("user_id")
	var userId uint
	if v, ok := userIdTemp.(uint); ok {
		userId = v
	}

	actionType := c.Query("action_type")
	videoIdStr := c.Query("video_id")
	videoId, _ := strconv.ParseUint(videoIdStr, 10, 64)

	// 发布评论
	if actionType == "1" {
		text := c.Query("comment_text")
		newComment, err1 := service.PostComment(userId, text, uint(videoId))
		commenter, err2 := service.GetCommenter(userId)
		author, err3 := service.GetAuthor(uint(videoId))

		if err1 != nil || err2 != nil || err3 != nil {
			c.JSON(http.StatusOK, common.BaseResponse{
				StatusCode: 403,
				StatusMsg:  "发表评论失败",
			})
			//c.Abort()
			log.Panicln("controller-Comment: 发表评论失败，")
			return
		}

		c.JSON(http.StatusOK, common.CommentActionBaseResp{
			BaseResponse: common.BaseResponse{
				StatusCode: 0,
				StatusMsg:  "发表评论成功",
			},
			Comment: common.CommentResp{
				ID:         newComment.ID,
				Content:    newComment.CommentText,
				CreateDate: newComment.CreatedAt.Format("01-02"),
				User:       util.PackUserInfo(commenter, service.CheckFollowing(userId, author)),
			},
		})
		return
	}

	if actionType != "2" {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 405,
			StatusMsg:  fmt.Sprint("未知的 actionType: ", actionType),
		})
		log.Panicf("controller-CommentAction: 评论操作失败：未知的 actionType：%v", actionType)
		return
	}

	//删除评论
	commentIdStr := c.Query("comment_id")
	commentId, _ := strconv.ParseInt(commentIdStr, 10, 64)
	err := service.DeleteComment(uint(commentId), uint(videoId))

	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 403,
			StatusMsg:  "Failed to delete comment",
		})
		c.Abort()
		log.Panicln("controller-Comment: 删除评论失败，")
		return
	}

	c.JSON(http.StatusOK, common.BaseResponse{
		StatusCode: 0,
		StatusMsg:  "删除评论成功",
	})

}

// CommentList 获取评论列表控制层
func CommentList(c *gin.Context) {

	token := c.Query("token")
	userID := util.GetHostIDFromToken(token)

	videoIdStr := c.Query("video_id")
	videoID, _ := strconv.ParseUint(videoIdStr, 10, 64)

	commentResponseList, err := service.CommentList(userID, uint(videoID))

	// log.Println("commentResponseList: ", commentResponseList)

	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 403,
			StatusMsg:  "获取评论列表失败",
		})
		c.Abort()
		log.Panicln("controller-CommentList: 获取评论列表失败，", err)
		return
	}

	//响应返回
	if len(commentResponseList) == 0 {
		c.JSON(http.StatusOK, common.CommentListBaseResp{
			BaseResponse: common.BaseResponse{
				StatusCode: 0,
				StatusMsg:  "该视频的评论列表为空",
			},
			CommentList: commentResponseList,
		})
		return
	}
	c.JSON(http.StatusOK, common.CommentListBaseResp{
		BaseResponse: common.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "获取评论列表成功",
		},
		CommentList: commentResponseList,
	})

}
