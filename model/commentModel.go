package model

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	UserID      int64
	VideoID     int64
	CommentText string
}
