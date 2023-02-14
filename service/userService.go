package service

import (
	"log"
	"strconv"
	"tiktok/mapper"
	"tiktok/pkg/common"
	"tiktok/pkg/errno"
	"tiktok/pkg/middleware"

	"golang.org/x/crypto/bcrypt"
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
	newUser, err := mapper.CreateUser(username, password)
	if err != nil {
		return userResponse, err
	}

	//3.颁发token
	token, err := middleware.CreateToken(newUser.ID, newUser.Name)
	if err != nil {
		log.Panicln("service-UserRegisterService: 创建用户token出错,", err)
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
		log.Panicln("service-UserInfoService: 解析rawID时发生错误， ", err)
		return userInfoQueryResponse, err
	}

	// 获取用户信息
	user, err := mapper.FindUserInfo(uint(userId))
	if err != nil {
		return userInfoQueryResponse, err
	}

	userInfoQueryResponse = common.UserInfoQueryResponse{
		UserID:        user.Model.ID,
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
		log.Println("用户名密码非法", err)
		return userResponse, err
	}

	// 检查用户是否存在
	user, flagExist := mapper.ExistUsername(username)
	if !flagExist {
		log.Println("Service-UserLoginService: 登录失败: 用户 ", username, " 不存在.")
		return userResponse, errno.ErrorFullPossibility
	}

	// 检查密码是否正确
	if !checkPassword(user.Password, password) {
		log.Println("service-UserLoginService: 登录失败：密码错误")
		return userResponse, errno.ErrorWrongPassword
	}

	if user.Model.ID == 0 {
		log.Println("账号或密码出错")
		return userResponse, errno.ErrorFullPossibility
	}

	// 颁发token
	token, err := middleware.CreateToken(user.Model.ID, user.Name)
	if err != nil {
		log.Println("service-UserLoginService: 登录失败，创建用户token发生错误")
		return userResponse, err
	}

	userResponse = common.UserIdTokenResponse{
		UserId: user.Model.ID,
		Token:  token,
	}

	return userResponse, nil
}

// IsFollow 检验已登录用户是否关注目标用户,
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
	return mapper.CheckFollowing(uint(hostId), userid)
}

// checkPassword 核对密码
func checkPassword(requestPassword string, truePassword string) bool {
	requestPasswordBytes := []byte(requestPassword)
	truePasswordBytes := []byte(truePassword)
	err := bcrypt.CompareHashAndPassword(requestPasswordBytes, truePasswordBytes)
	if err != nil {
		log.Panicln("userService-checkPassword: 核对密码时出错，", err)
	}
	return true
}

// isLegal 检查用户名和密码的合法性
func isLegal(username string, password string) error {
	//1.用户名检验
	if username == "" {
		log.Println("用户名为空")
		return errno.ErrorNullUsername
	}
	if len(username) > MaxUsernameLength {
		log.Printf("用户名过长，应小于%d位\n", MaxUsernameLength)
		return errno.ErrorUsernameExtend
	}

	//2.密码检验
	if password == "" {
		log.Panicln("密码为空")
		return errno.ErrorNullPassword
	}
	if len(password) > MaxPasswordLength {
		log.Printf("密码过长，应小于%d位\n", MaxPasswordLength)
		return errno.ErrorPasswordLength
	}
	return nil
}
