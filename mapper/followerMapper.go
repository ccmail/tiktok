package mapper

import (
	"errors"
	"tiktok/model"

	"gorm.io/gorm"
)

// IsFollowing 判断 FollowerID 是否关注 UserID
func CheckFollowing(FollowerID uint, UserID uint) bool {
	var relationExist = &model.Follower{}
	//判断关注是否存在
	err := DBConn.Model(&model.Follower{}).Where("follower_id=? AND user_id=? AND is_follow=?", FollowerID, UserID, true).First(&relationExist).Error

	// false-关注不存在，true-关注存在
	return !errors.Is(err, gorm.ErrRecordNotFound)
}
