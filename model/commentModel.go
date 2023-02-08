package model

import (
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	//默认设置为主键
	ID uint64 `gorm:"primaryKey"`
	//gorm定义的model, 包含id , 创建时间, 更新时间, 删除时间
	UserID      int64
	VideoID     int64
	CommentText string
	//gorm自动填充创建, 更改, 删除时间
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
