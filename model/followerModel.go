package model

import "gorm.io/gorm"

type Follower struct {
	gorm.Model
	//粉丝id
	FollowerID uint
	//up id
	UserID uint
	//冗余字段, 避免多表查询带来的IO损失
	FollowerName string
	//是否关注user, 默认值设为false, false不关注, true关注
	IsFollow bool `gorm:"default:false"`
}
