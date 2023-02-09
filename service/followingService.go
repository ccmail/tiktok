package service

import (
	"errors"
	"tiktok/mapper"
	"tiktok/model"

	"gorm.io/gorm"
)

// IsFollowing 判断 FollowerID 是否关注 UserID
func IsFollowing(FollowerID uint, UserID uint) bool {
	var relationExist = &model.Follower{}
	//判断关注是否存在
	err := mapper.DBConn.Model(&model.Follower{}).Where("FollowerID=? AND UserID=? AND IsFollow=?", FollowerID, UserID, true).First(&relationExist).Error

	// false-关注不存在，true-关注存在
	return errors.Is(err, gorm.ErrRecordNotFound)
}
