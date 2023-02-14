package controller

import (
	"log"
	"net/http"
	"strconv"
	"tiktok/pkg/common"
	"tiktok/pkg/middleware"
	"tiktok/service"

	"github.com/gin-gonic/gin"
)

// CommentAction 评论操作控制层
func Comment(c *gin.Context) {
	//1 数据处理
	UserId, _ := c.Get("user_id")
	var userId uint
	if v, ok := UserId.(uint); ok {
		userId = v
	}
	actionType := c.Query("action_type")
	videoIdStr := c.Query("video_id")
	videoId, _ := strconv.ParseUint(videoIdStr, 10, 64)

	if actionType == "1" { // 发布评论
		text := c.Query("comment_text")

		// TODO: 这个调用的3个service，其实可以并发执行，后面可以考虑用协程来实现并发
		newComment, err1 := service.PostCommentService(userId, text, uint(videoId))
		commenter, err2 := service.GetCommenter(userId)
		author, err3 := service.GetAuthor(uint(videoId))
		if err1 != nil || err2 != nil || err3 != nil {
			c.JSON(http.StatusOK, common.BaseResponse{
				StatusCode: 403,
				StatusMsg:  "发表评论失败",
			})
			c.Abort()
			log.Panicln("controller-Comment: 发表评论失败，")
		}

		c.JSON(http.StatusOK, common.CommentActionResponse{
			BaseResponse: common.BaseResponse{
				StatusCode: 0,
				StatusMsg:  "发表评论成功",
			},
			Comment: common.CommentResponse{
				ID:         newComment.ID,
				Content:    newComment.CommentText,
				CreateDate: newComment.CreatedAt.Format("01-02"),
				User: common.CommenterInfo{
					UserID:        commenter.ID,
					Username:      commenter.Name,
					FollowCount:   commenter.FollowCount,
					FollowerCount: commenter.FollowerCount,
					// TODO: IsFollowing包含了DAO操作，后面要移到mapper里去
					IsFollow: service.IsFollowing(userId, author),
				},
			},
		})

	} else if actionType == "2" { //删除评论
		commentIdStr := c.Query("comment_id")
		commentId, _ := strconv.ParseInt(commentIdStr, 10, 64)
		err := service.DeleteCommentService(uint(videoId), uint(commentId))

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

	} else {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 405,
			StatusMsg:  "未知的 actionType: " + actionType,
		})
		c.Abort()
		log.Println("controller-CommentAction: 评论操作失败：未知的 actionType：" + actionType)
		return
	}

}

func getHostIDFromToken(hostToken string) uint {
	var hostID uint
	if hostToken != "" {
		hostInfo, err := middleware.ParseTokenCJS(hostToken)
		if err != nil {
			//就算token解析失败也不应该拒绝访问publishList
			log.Println("请求评论列表的用户的token解析失败", err)
		} else {
			hostID = hostInfo.UserId
		}
	}
	//log.Println(hostID)
	return hostID
}

// CommentList 获取评论列表控制层
func CommentList(c *gin.Context) {
	// getUserID, _ := c.Get("user_id")
	// var userID uint
	// if v, ok := getUserID.(uint); ok {
	// 	userID = v
	// }
	token := c.Query("token")
	userID := getHostIDFromToken(token)

	videoIdStr := c.Query("video_id")
	videoID, _ := strconv.ParseUint(videoIdStr, 10, 64)

	// log.Println("userID: ", userID, ", videoID: ", videoID)

	commentRespionseList, err := service.CommentListService(userID, uint(videoID))

	// log.Println("commentRespionseList: ", commentRespionseList)

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
	if len(commentRespionseList) == 0 {
		c.JSON(http.StatusOK, common.CommentListResponse{
			BaseResponse: common.BaseResponse{
				StatusCode: 0,
				StatusMsg:  "该视频的评论列表为空",
			},
			CommentList: commentRespionseList,
		})
		return
	}
	c.JSON(http.StatusOK, common.CommentListResponse{
		BaseResponse: common.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "获取评论列表成功",
		},
		CommentList: commentRespionseList,
	})

}
