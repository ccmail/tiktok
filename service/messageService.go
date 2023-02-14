package service

import (
	"log"
	"tiktok/mapper"
	"tiktok/model"
	"tiktok/pkg/common"
)

// SendMessageService 发送消息服务
func SendMessageService(senderID uint, receiverID uint, message_text string) error {
	message := model.Message{
		UserID:      senderID,
		FriendID:    receiverID,
		MessageText: message_text,
	}
	err := mapper.CreateMessage(message)
	if err != nil {
		log.Println("service-SendMessageService: 发送消息失败，", err.Error())
		return err
	}
	return nil
}

func GetMessageListService(senderID uint, receiverID uint) (messageResponseList []common.MessageResponse, err error) {
	messageList, err := mapper.GetMessageList(senderID, receiverID)
	if err != nil {
		log.Println("service-SendMessageService: 获取消息列表失败，", err.Error())
		return []common.MessageResponse{}, err
	}

	for _, m := range messageList {
		var messageResponse = common.MessageResponse{
			ID:         m.ID,
			Content:    m.MessageText,
			CreateTime: m.CreatedAt.Format("2006-01-02 15:03:04"),
		}
		messageResponseList = append(messageResponseList, messageResponse)
	}
	return messageResponseList, nil
}
