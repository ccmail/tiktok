package service

import (
	"log"
	"strconv"
	"tiktok/mapper/db"
	"tiktok/model"
	"tiktok/pkg/common"
	"time"
)

// SendMessage 发送消息服务
func SendMessage(senderID uint, receiverID uint, messageText string) error {
	message := model.Message{
		UserID:      senderID,
		FriendID:    receiverID,
		MessageText: messageText,
	}
	err := db.CreateMessage(message)
	if err != nil {
		log.Println("service-SendMessage: 发送消息失败，", err.Error())
		return err
	}
	return nil
}

func GetMessageList(senderID uint, receiverID uint, strPrevTime string) (messageResponseList []common.MessageResp, err error) {
	prevTime := time.Unix(0, 0)
	if strPrevTime != "0" && strPrevTime != "" {
		i, err := strconv.ParseInt(strPrevTime, 10, 64)
		if err != nil {
			log.Println("传入的时间有问题")
		} else {
			//转化为本地时间
			prevTime = time.UnixMilli(i).Local()
		}
	}
	messageList, err := db.GetMessageList(senderID, receiverID, prevTime)
	log.Println("打印时间", prevTime)
	if err != nil {
		log.Println("service-SendMessage: 获取消息列表失败，", err.Error())
		return []common.MessageResp{}, err
	}

	for _, m := range messageList {
		var messageResponse = common.MessageResp{
			ID:         m.ID,
			Content:    m.MessageText,
			FromUserID: m.UserID,
			ToUserID:   m.FriendID,
			//CreateTime: m.CreatedAt.Format("2006-01-02 15:04:05.000"),
			CreateTime: m.CreatedAt.UnixMilli(),
		}
		messageResponseList = append(messageResponseList, messageResponse)
	}
	return messageResponseList, nil
}
