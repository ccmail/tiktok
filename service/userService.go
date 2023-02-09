package service

import (
	"errors"
	"strconv"
	"tiktok/mapper"
	"tiktok/pkg/common"
	"tiktok/pkg/errno"
	middleware "tiktok/pkg/mw"

	"tiktok/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	MaxUsernameLength = 32 //用户名最大长度
	MaxPasswordLength = 32 //密码最大长度
	MinPasswordLength = 8  //密码最小长度
)

// UserRegisterService 用户注册服务
func UserRegisterService(username string, password string) (common.UserIdTokenResponse, error) {

	//0.数据准备
	var userResponse = common.UserIdTokenResponse{}

	//1.合法性检验
	err := isLegal(username, password)
	if err != nil {
		return userResponse, err
	}

	//2.新建用户
	newUser, err := createUser(username, password)
	if err != nil {
		return userResponse, err
	}

	//3.颁发token
	token, err := middleware.CreateToken(newUser.ID, newUser.Name)
	if err != nil {
		return userResponse, err
	}

	userResponse = common.UserIdTokenResponse{
		UserId: newUser.ID,
		Token:  token,
	}
	return userResponse, nil
}

// UserInfoService 用户信息获取服务
func UserInfoService(rawId string) (common.UserInfoQueryResponse, error) {
	// 数据准备
	var userInfoQueryResponse = common.UserInfoQueryResponse{}
	userId, err := strconv.ParseUint(rawId, 10, 64)
	if err != nil {
		return userInfoQueryResponse, err
	}

	// 获取用户信息
	user, err := GetUser(uint(userId))
	if err != nil {
		return userInfoQueryResponse, err
	}

	userInfoQueryResponse = common.UserInfoQueryResponse{
		UserId:        user.Model.ID,
		Username:      user.Name,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      false,
	}
	return userInfoQueryResponse, nil
}

// UserLoginService 用户登录服务
func UserLoginService(username string, password string) (common.UserIdTokenResponse, error) {

	// 数据准备
	var userResponse = common.UserIdTokenResponse{}

	// 合法性检验
	err := isLegal(username, password)
	if err != nil {
		return userResponse, err
	}

	// 检查用户是否存在
	var user model.User

	err = mapper.DBConn.Where("name=?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return userResponse, errno.ErrorFullPossibility
	}

	// 检查密码是否正确
	if !checkPassword(user.Password, password) {
		return userResponse, errno.ErrorWrongPassword
	}

	if user.Model.ID == 0 {
		return userResponse, errno.ErrorFullPossibility
	}

	// 颁发token
	token, err := middleware.CreateToken(user.Model.ID, user.Name)
	if err != nil {
		return userResponse, err
	}

	userResponse = common.UserIdTokenResponse{
		UserId: user.Model.ID,
		Token:  token,
	}

	return userResponse, nil
}

// GetUser 根据用户id获取用户信息
func GetUser(userId uint) (model.User, error) {
	//1.数据模型准备
	var user model.User
	//2.在users表中查对应user_id的user
	err := mapper.DBConn.Model(&model.User{}).Where("id = ?", userId).Find(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

// IsFollow 检验已登录用户是否关注目标用户
func IsFollow(targetId string, userid uint) bool {
	// 修改targetId数据类型
	hostId, err := strconv.ParseUint(targetId, 10, 64)
	if err != nil {
		return false
	}
	// 如果是自己查自己，那就是没有关注
	if uint(hostId) == userid {
		return false
	}
	// 自己是否关注目标userId
	return IsFollowing(uint(hostId), userid)
}

// createUser 新建用户
func createUser(username string, password string) (model.User, error) {
	// Following数据模型准备
	encryptedPassword, _ := encrypt(password)
	newUser := model.User{
		Name:     username,
		Password: encryptedPassword,
	}
	// 模型关联到数据库表users //可注释
	mapper.DBConn.AutoMigrate(&model.User{})
	// 新建user
	if existUsername(username) {
		//用户名已存在
		return newUser, errno.ErrorRedundantUsername
	}

	// 用户不存在，在DB中新建用户
	err := mapper.DBConn.Model(&model.User{}).Create(&newUser).Error
	if err != nil {
		// 错误处理
		panic(err)
	}

	return newUser, nil
}

// existUsername 检查用户名是否存在
func existUsername(username string) bool {
	var userExist = &model.User{}
	err := mapper.DBConn.Model(&model.User{}).Where("name=?", username).First(&userExist).Error

	// false-用户名不存在，true-用户名存在
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// checkPassword 核对密码
func checkPassword(requestPassword string, truePassword string) bool {
	requestPasswordBytes := []byte(requestPassword)
	truePasswordBytes := []byte(truePassword)
	err := bcrypt.CompareHashAndPassword(requestPasswordBytes, truePasswordBytes)
	if err != nil {
		panic(err)
	}
	return true
}

// encrypt 使用 bcrypt 对密码进行加密
func encrypt(passwordString string) (encryptedPassword string, err error) {
	passwdBytes := []byte(passwordString)
	hash, err := bcrypt.GenerateFromPassword(passwdBytes, bcrypt.MinCost)
	if err != nil {
		return
	}
	encryptedPassword = string(hash)
	return
}

// isLegal 检查用户名和密码的合法性
func isLegal(username string, password string) error {
	//1.用户名检验
	if username == "" {
		return errno.ErrorUserNameNull
	}
	if len(username) > MaxUsernameLength {
		return errno.ErrorUserNameExtend
	}

	//2.密码检验
	if password == "" {
		return errno.ErrorPasswordNull
	}
	if len(password) > MaxPasswordLength {
		return errno.ErrorPasswordLength
	}
	return nil
}
