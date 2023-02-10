package model

import "gorm.io/gorm"

// Like 用户点赞的作品列表
type Like struct {
	gorm.Model
	UserID  uint
	VideoID uint64
	// 视频封面
	CoverUrl string
	VideoURL string
	// 用户是否对这个作品点赞, 默认为true, 当用户取消点赞时, 将这一条IsLike设置为False
	IsLike bool `gorm:"default:true"`
}
