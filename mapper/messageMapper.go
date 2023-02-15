package mapper

import "tiktok/model"

// CreateMessage 创建一条消息
func CreateMessage(message model.Message) error {
	create := DBConn.Table("messages").Create(&message)
	if create.Error != nil {
		return create.Error
	}
	return nil
}

// GetMessageList 获取消息列表
func GetMessageList(senderID uint, receiverID uint) (messageList []model.Message, err error) {
	find := DBConn.Table("messages").Where("user_id = ? AND friend_id = ?", senderID, receiverID).Find(&messageList)
	if find.Error != nil {
		return []model.Message{}, find.Error
	}
	return messageList, nil
}
