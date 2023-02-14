package mapper

import (
	"errors"
	"log"
	"tiktok/model"

	"gorm.io/gorm"
)

// CheckFollowing  判断 FollowerID 是否关注 UserID
func CheckFollowing(FollowerID uint, UserID uint) bool {
	if FollowerID == 0 {
		//用户处于未登录的状态, 默认未关注
		return false
	}
	var relationExist = &model.Follower{}
	//判断关注是否存在
	err := DBConn.Model(&model.Follower{}).Where("follower_id=? AND user_id=? AND is_follow=?", FollowerID, UserID, true).First(&relationExist).Error

	// false-关注不存在，true-关注存在
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// ExistFollowRecord
// 是host去关注guest, 这里host是粉丝, guestID是up
func ExistFollowRecord(hostID, guestID uint) (followRecord model.Follower, exist bool) {
	find := DBConn.Model(&model.Follower{}).Where("user_id = ? AND follower_id = ?", guestID, hostID).Limit(1).Find(&followRecord)
	if find.Error != nil {
		log.Panicln("查找失败")
		return followRecord, false
	}
	if followRecord.ID == 0 {
		log.Println("没有查找到相关关注记录")
		return followRecord, false
	}
	return followRecord, true
}

// CreatFollowRecord
// host关注guest, 所以host是被关注的人, guest是up
func CreatFollowRecord(hostID, guestID uint, isConcern bool) error {
	followRecord := model.Follower{IsFollow: isConcern, FollowerID: hostID, UserID: guestID}
	//正在关注, 执行+1
	tx := DBConn.Model(&model.Follower{}).Create(&followRecord)
	if tx.Error != nil {
		log.Panicln("插入关注信息时出错")
		return tx.Error
	}
	return nil
}

// UpdateFollowRecord 是host去关注guest, 这里host是follower
func UpdateFollowRecord(hostID, guestID uint, isConcern bool) error {
	tx := DBConn.Model(&model.Follower{}).Where("user_id = ? AND follower_id = ?", guestID, hostID).Update("IsFollow", isConcern)
	if tx.Error != nil {
		log.Panicln("更新关注信息时出错")
		return tx.Error
	}
	return nil
}

// FindMultiConcern 返回关注id关注的人, 实现逻辑是将id作为粉丝id进行查询
func FindMultiConcern(id uint) (resUserIDList []uint, err error) {
	tx := DBConn.Model(&model.Follower{}).Select("user_id").Where("follower_id = ? AND is_follow = ?", id, true).Find(&resUserIDList)
	if tx.Error != nil {
		log.Panicln("查询用户关注信息时失败")
		return resUserIDList, tx.Error
	}
	return resUserIDList, nil
}

func FindMultiFollower(id uint) (resUserIDList []uint, err error) {
	tx := DBConn.Model(&model.Follower{}).Select("follower_id").Where("user_id = ? AND is_follow = ?", id, true).Find(&resUserIDList)
	if tx.Error != nil {
		log.Panicln("查询用户关注信息时失败")
		return resUserIDList, tx.Error
	}
	return resUserIDList, nil
}
