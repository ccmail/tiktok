package service

import (
	"tiktok/mapper"
)

// IsFollowing 判断 FollowerID 是否关注 UserID
func IsFollowing(FollowerID uint, UserID uint) bool {
	if FollowerID == 0 {
		//用户处于未登录的状态, 默认未关注
		return false
	}
	return mapper.CheckFollowing(FollowerID, UserID)
}
