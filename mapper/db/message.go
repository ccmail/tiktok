package db

import (
	"tiktok/mapper"
	"tiktok/model"
	"time"
)

// CreateMessage 创建一条消息
func CreateMessage(message model.Message) error {
	create := mapper.DBConn.Table("messages").Create(&message)
	if create.Error != nil {
		return create.Error
	}
	return nil
}

// GetMessageList 获取消息列表
func GetMessageList(senderID uint, receiverID uint, prevTime time.Time) (messageList []model.Message, err error) {
	find := mapper.DBConn.Table("messages").
		Where("user_id = ? AND friend_id = ? OR user_id = ? AND friend_id = ? ", senderID, receiverID, receiverID, senderID).
		Where("created_at > ?", prevTime).
		Find(&messageList)
	if find.Error != nil {
		return []model.Message{}, find.Error
	}
	return messageList, nil
}

func GetSendMessage(hostID, guestID uint) (message model.Message) {
	last := mapper.DBConn.Model(&model.Message{}).Where("user_id = ? AND friend_id = ?", hostID, guestID).Last(&message)
	if last.Error != nil {
		//	没有查询到对应的最后一条信息, 因此返回""
		return
	}
	return message
}

func GetReceiveMessage(hostID, guestID uint) (message model.Message) {
	last := mapper.DBConn.Model(&model.Message{}).Where("user_id = ? AND friend_id = ?", guestID, hostID).Last(&message)
	if last.Error != nil {
		return
	}
	return message
}
