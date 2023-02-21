package model

import "gorm.io/gorm"

// User 暂时的字段,后期按照提供的app详情可能会需要完善字段
type User struct {
	gorm.Model           //db.Model里包含了ID字段且默认为主键
	Name          string `db:"unique;not null;"` // 用户名称
	Password      string `db:"unique;not null;"` // 密码
	FollowCount   uint
	FollowerCount uint
}
