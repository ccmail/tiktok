package mapper

import (
	"log"
	"tiktok/model"
)

func FindUserInfo(userId uint) (user model.User, err error) {
	err = DBConn.Model(&model.User{}).Where("id = ?", userId).Find(&user).Error
	if err != nil {
		log.Panicf("没有查到id为%v的用户\n", userId)
		return user, err
	}
	return user, nil
}

// FindMultiUserInfo 返回值为map, 根据userID获取user结构体
func FindMultiUserInfo(multiUserID []uint) (map[uint]model.User, error) {
	//根据userID为索引, 存储user信息
	mp := map[uint]model.User{}
	//获取查询到的user信息
	var tempUser []model.User
	find := DBConn.Model(&model.User{}).Where("id IN ?", multiUserID).Find(&tempUser)
	if find.Error != nil {
		log.Panicln("查询多个userInfo时发生了错误", find.Error)
		return mp, find.Error
	}
	//将获取到的user信息存到map中
	for _, userInfo := range tempUser {
		mp[userInfo.ID] = userInfo
	}
	return mp, nil
}
