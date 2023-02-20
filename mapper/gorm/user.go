package gorm

import (
	"errors"
	"fmt"
	"log"
	"tiktok/config"
	"tiktok/model"
	"tiktok/pkg/errno"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// FindUserInfo 根据用户id获取用户信息
func FindUserInfo(userId uint) (user model.User, err error) {
	err = config.DBConn.Model(&model.User{}).Where("id = ?", userId).Find(&user).Error
	if err != nil {
		log.Printf("没有查到id为%v的用户\n", userId)
		return user, nil
	}
	return user, nil
}

// CreateUser  新建用户
func CreateUser(username string, password string) (model.User, error) {
	// Following数据模型准备
	encryptedPassword, _ := encrypt(password)
	newUser := model.User{
		Name:     username,
		Password: encryptedPassword,
	}
	// 模型关联到数据库表users //可注释
	err := config.DBConn.AutoMigrate(&model.User{})
	if err != nil {
		log.Panicln("模型关联表失败")
		return newUser, err
	}
	// 新建user
	if _, flagExist := ExistUsername(username); flagExist {
		//用户名已存在
		log.Println("mapper-CreateUser: 无法创建用户：用户名已存在")
		return newUser, errno.ErrorRedundantUsername
	}

	// 用户不存在，在DB中新建用户
	err = config.DBConn.Model(&model.User{}).Create(&newUser).Error
	if err != nil {
		// 错误处理
		log.Panicln("mapper-CreateUser: 创建用户时出错", err)
	}

	return newUser, nil
}

// encrypt 使用 bcrypt 对密码进行加密
func encrypt(passwordString string) (encryptedPassword string, err error) {
	passwdBytes := []byte(passwordString)
	hash, err := bcrypt.GenerateFromPassword(passwdBytes, bcrypt.MinCost)
	if err != nil {
		log.Panicln("加密密码失败, ", err)
		return
	}
	encryptedPassword = string(hash)
	return
}

// ExistUsername  检查用户名是否存在
func ExistUsername(username string) (model.User, bool) {
	var user model.User
	err := config.DBConn.Model(&model.User{}).Where("name=?", username).First(&user).Error

	// false-用户名不存在，true-用户名存在
	return user, !errors.Is(err, gorm.ErrRecordNotFound)
}

// FindMultiUserInfo 返回值为map, 根据userID获取user结构体
func FindMultiUserInfo(multiUserID []uint) (map[uint]model.User, error) {
	//根据userID为索引, 存储user信息
	mp := map[uint]model.User{}
	//获取查询到的user信息
	var tempUser []model.User
	find := config.DBConn.Model(&model.User{}).Where("id IN ?", multiUserID).Find(&tempUser)
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

// UpdateUserFollowCount host关注guest, 所以host是被关注的人, guest是up, 所以关注操作, host.follow+1, guest.follower+1
func UpdateUserFollowCount(hostID, guestID uint, isConcern bool) error {
	x := " + 1"
	if !isConcern {
		x = " - 1"
	}
	tx := config.DBConn.Model(&model.User{}).Where("id = ?", hostID).Update("follow_count", gorm.Expr(fmt.Sprint("follow_count", x)))
	if tx.Error != nil {
		log.Panicln("更新关注人数时出错")
		return tx.Error
	}
	tx = config.DBConn.Model(&model.User{}).Where("id = ?", guestID).Update("follower_count", gorm.Expr(fmt.Sprint("follower_count ", x)))
	if tx.Error != nil {
		log.Panicln("更新粉丝人数时出错")
		return tx.Error
	}
	return nil
}

func GetMultiUserInfoNoHit(userInfo *[]model.User, userNoCache *map[uint][]int) (err error) {
	users := make([]uint, 0, len(*userNoCache))
	for k := range *userNoCache {
		users = append(users, k)
	}

	var userList []model.User
	find := config.DBConn.Model(&model.User{}).Where("id IN ?", users).Find(&userList)
	if find.Error != nil {
		log.Panicln("在mysql查询用户信息失败")
		return err
	}

	for _, user := range userList {
		for _, v := range (*userNoCache)[user.ID] {
			(*userInfo)[v] = user
		}
	}
	return nil
}
