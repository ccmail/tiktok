package mapper

import (
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	//ID默认设置为主键, 指明uint64, 避免在一些32位机器上出错
	ID uint64 `gorm:"primaryKey"`
	//gorm定义的model, 包含id , 创建时间, 更新时间, 删除时间
	UserID      uint64
	VideoID     uint64
	CommentText string
	//gorm自动填充创建, 更改, 删除时间
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (Comment) TableName() string {
	//默认返回也是蛇形复数, 增加熟练度
	return "comments"
}

func Count() {

}
