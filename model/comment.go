package model

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	UserID      uint
	VideoID     uint `gorm:"index:vid"`
	CommentText string
	Valid       bool `db:"default:true"`
}
