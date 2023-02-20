package model

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	UserID      uint
	VideoID     uint
	CommentText string
	Valid       bool `gorm:"default:true"`
}
