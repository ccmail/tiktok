package mapper

import (
	"log"
	"tiktok/model"
)

func FindUserInfo(userId uint) (user model.User, err error) {
	err = DBConn.Model(&model.User{}).Where("id = ?", userId).Find(&user).Error
	if err != nil {
		log.Printf("没有查到id为%v的用户\n", userId)
		return user, nil
	}
	return user, nil
}
