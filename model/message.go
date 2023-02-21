package model

import "gorm.io/gorm"

type Message struct {
	//消息的ID
	gorm.Model
	UserID      uint `gorm:"index:user_friend_union_idx"`
	FriendID    uint `gorm:"index:user_friend_union_idx"`
	MessageText string
}
