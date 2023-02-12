package mapper

import (
	"errors"
	"log"
	"tiktok/model"
	"tiktok/pkg/errno"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// FindUserInfo 根据用户id获取用户信息
func FindUserInfo(userId uint) (user model.User, err error) {
	err = DBConn.Model(&model.User{}).Where("id = ?", userId).Find(&user).Error
	if err != nil {
		log.Printf("没有查到id为%v的用户\n", userId)
		return user, nil
	}
	return user, nil
}

// createUser 新建用户
func CreateUser(username string, password string) (model.User, error) {
	// Following数据模型准备
	encryptedPassword, _ := encrypt(password)
	newUser := model.User{
		Name:     username,
		Password: encryptedPassword,
	}
	// 模型关联到数据库表users //可注释
	DBConn.AutoMigrate(&model.User{})
	// 新建user
	if _, flagExist := ExistUsername(username); flagExist {
		//用户名已存在
		log.Println("mapper-CreateUser: 无法创建用户：用户名已存在")
		return newUser, errno.ErrorRedundantUsername
	}

	// 用户不存在，在DB中新建用户
	err := DBConn.Model(&model.User{}).Create(&newUser).Error
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

// existUsername 检查用户名是否存在
func ExistUsername(username string) (model.User, bool) {
	var user model.User
	err := DBConn.Model(&model.User{}).Where("name=?", username).First(&user).Error

	// false-用户名不存在，true-用户名存在
	return user, !errors.Is(err, gorm.ErrRecordNotFound)
}
