package model

import "gorm.io/gorm"

type Follower struct {
	gorm.Model
	FollowerID uint64
	UserID     uint64
	//冗余字段, 避免多表查询带来的IO损失
	FollowerName string
	//是否不再关注user, 默认值设为false,
	IsFollow bool `gorm:"default:false"`
}
