package controller

import (
	"log"
	"net/http"
	"strconv"
	"tiktok/pkg/common"
	"tiktok/service"

	"github.com/gin-gonic/gin"
)

// Message 发送消息控制层
func Message(c *gin.Context) {
	getUserId, _ := c.Get("user_id")
	var senderID uint
	if v, ok := getUserId.(uint); ok {
		senderID = v
	}
	receiverIDStr := c.Query("to_user_id")
	receiverID, err := strconv.ParseUint(receiverIDStr, 10, 64)
	if err != nil || receiverIDStr == "" {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "发送消息失败，" + err.Error(),
		})
		log.Panicln("controller-Message: 发送消息失败，解析接收方ID时出错，" + err.Error())
	}

	actionTypeStr := c.Query("action_type")
	actionType, err := strconv.ParseUint(actionTypeStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "发送消息失败，" + err.Error(),
		})
		log.Panicln("controller-Message: 发送消息失败，解析actionType时出错，" + err.Error())
	}

	message_text := c.Query("content")

	if actionType != 1 {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "发送消息失败：未知的 actionType: " + actionTypeStr,
		})
		log.Panicln("controller-Message: 发送消息失败，未知的 actionType: " + actionTypeStr)
	}

	err = service.SendMessageService(senderID, uint(receiverID), message_text)
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "controller-Message: 发送消息失败，" + err.Error(),
		})
		log.Panicln("controller-Message: 发送消息失败，" + err.Error())
	}
	c.JSON(http.StatusOK, common.BaseResponse{
		StatusCode: 0,
		StatusMsg:  "发送消息成功",
	})

}

// 消息记录控制层
func MessageList(c *gin.Context) {
	getUserId, _ := c.Get("user_id")
	var senderID uint
	if v, ok := getUserId.(uint); ok {
		senderID = v
	}
	receiverIDStr := c.Query("to_user_id")
	receiverID, err := strconv.ParseUint(receiverIDStr, 10, 64)
	if err != nil || receiverIDStr == "" {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "获取消息记录失败，" + err.Error(),
		})
		log.Panicln("controller-Message: 获取消息记录失败，解析接收方ID时出错，" + err.Error())
	}
	messageResponseList, err := service.GetMessageListService(senderID, uint(receiverID))

	if err != nil || receiverIDStr == "" {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "获取消息记录失败，" + err.Error(),
		})
		log.Panicln("controller-Message: 获取消息记录失败，" + err.Error())
	}

	if len(messageResponseList) == 0 {
		c.JSON(http.StatusOK, common.MessageListResponse{
			BaseResponse: common.BaseResponse{
				StatusCode: 0,
				StatusMsg:  "获取消息记录成功，消息历史为空",
			},
			MessageResponseList: messageResponseList,
		})
		return
	}

	c.JSON(http.StatusOK, common.MessageListResponse{
		BaseResponse: common.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "获取消息记录成功",
		},
		MessageResponseList: messageResponseList,
	})

}
