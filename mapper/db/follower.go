package db

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"tiktok/mapper"
	"tiktok/model"
)

// CheckFollowing  判断 hostID 是否关注 guestID
func CheckFollowing(hostID uint, guestID uint) bool {
	if hostID == 0 {
		//用户处于未登录的状态, 默认未关注
		return false
	}
	var relationExist = &model.Follower{}
	//判断关注是否存在
	err := mapper.DBConn.Model(&model.Follower{}).Where("follower_id=? AND user_id=? AND is_follow=?", hostID, guestID, true).First(&relationExist).Error

	// false-关注不存在，true-关注存在
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

func TranCreateFollow(hostID, guestID uint, isConcern bool) (err error) {
	err = mapper.DBConn.Transaction(func(tx *gorm.DB) error {
		if err = CreatFollowRecord(hostID, guestID, isConcern); err != nil {
			log.Panicln("插入数据库时发生错误")
			return err
		}
		//从未关注过的话直接插, 更新完毕之后return掉
		if err = UpdateUserFollowCount(hostID, guestID, isConcern); err != nil {
			log.Panicln("更新粉丝/关注数量时出错")
			return err
		}
		return nil
	})
	if err != nil {
		log.Panicln("关注操作失败!", err)
	}
	return err
}

func TranUpdateFollow(hostID, guestID uint, isConcern bool) (err error) {
	err = mapper.DBConn.Transaction(func(tx *gorm.DB) error {

		err = UpdateUserFollowCount(hostID, guestID, isConcern)
		if err != nil {
			log.Panicln("更新粉丝/关注数量时出错")
			return err
		}
		//向follow表中添加follow关系记录
		err = UpdateFollowRecord(hostID, guestID, isConcern)
		if err != nil {
			log.Panicln("更新关注关系时失败")
			return err
		}
		return nil
	})
	if err != nil {
		log.Panicln("更新关注关系失败")
	}
	return nil
}

// CheckFollowRecord
// 是host去关注guest, 这里host是粉丝, guestID是up
func CheckFollowRecord(hostID, guestID uint) (followRecord model.Follower, exist bool) {
	find := mapper.DBConn.Model(&model.Follower{}).Where("user_id = ? AND follower_id = ?", guestID, hostID).Limit(1).Find(&followRecord)
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
// host关注guest, 所以host是粉丝, guest是up
func CreatFollowRecord(hostID, guestID uint, isConcern bool) error {
	followRecord := model.Follower{IsFollow: isConcern, FollowerID: hostID, UserID: guestID}
	//正在关注, 执行+1
	tx := mapper.DBConn.Model(&model.Follower{}).Create(&followRecord)
	if tx.Error != nil {
		log.Panicln("插入关注信息时出错")
		return tx.Error
	}
	return nil
}

// UpdateFollowRecord 是host去关注guest, 这里host是follower
func UpdateFollowRecord(hostID, guestID uint, isConcern bool) error {
	tx := mapper.DBConn.Model(&model.Follower{}).Where("user_id = ? AND follower_id = ?", guestID, hostID).Update("IsFollow", isConcern)
	if tx.Error != nil {
		log.Panicln("更新关注信息时出错")
		return tx.Error
	}
	return nil
}

// GetMultiConcern 返回关注id关注的人, 实现逻辑是将id作为粉丝id进行查询
func GetMultiConcern(id uint) (resUserIDList []uint, err error) {
	tx := mapper.DBConn.Model(&model.Follower{}).Select("user_id").Where("follower_id = ? AND is_follow = ?", id, true).Find(&resUserIDList)
	if tx.Error != nil {
		log.Panicln("查询用户关注信息时失败")
		return resUserIDList, tx.Error
	}
	return resUserIDList, nil
}

func FindMultiFollower(id uint) (resUserIDList []uint, err error) {
	tx := mapper.DBConn.Model(&model.Follower{}).Select("follower_id").Where("user_id = ? AND is_follow = ?", id, true).Find(&resUserIDList)
	if tx.Error != nil {
		log.Panicln("查询用户关注信息时失败")
		return resUserIDList, tx.Error
	}
	return resUserIDList, nil
}

// CheckMultiFollowNoHit 判断hostID是否关注followNoCache中的人
func CheckMultiFollowNoHit(hostID uint, isFollow *[]bool, followNoCache *map[uint][]int) (err error) {
	follows := make([]uint, 0, len(*followNoCache)>>2)
	for k := range *followNoCache {
		follows = append(follows, k)
	}
	var followList []model.Follower
	err = mapper.DBConn.Model(&model.Follower{}).Where("follower_id = ? AND  user_id IN ?", hostID, follows).Find(&followList).Error
	if err != nil {
		log.Println("在mysql查询未命中的关注关系时失败")
		return err
	}
	for _, f := range followList {
		for _, v := range (*followNoCache)[f.UserID] {
			(*isFollow)[v] = f.IsFollow
		}
	}
	return nil
}
