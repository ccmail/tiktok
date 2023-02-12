package mapper

import (
	"errors"

	"gorm.io/gorm"
)

// IsFavorite 查询某用户是否点赞某视频
func IsFavorite(uid uint, vid uint) bool {
	var total int64
	err := DBConn.Table("likes").Where("user_id = ? AND video_id = ? AND is_like = ?", uid, vid, true).Count(&total).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	if total == 0 {
		return false
	}
	return true
}
