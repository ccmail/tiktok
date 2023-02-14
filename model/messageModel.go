package model

import "gorm.io/gorm"

type Message struct {
	//消息的ID
	gorm.Model
	UserID      uint
	FriendID    uint
	MessageText string
}
