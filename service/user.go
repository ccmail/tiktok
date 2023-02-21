package service

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"strconv"
	"tiktok/mapper/cache"
	"tiktok/mapper/db"
	"tiktok/pkg/common"
	"tiktok/pkg/errno"
	"tiktok/pkg/middleware"
	"tiktok/pkg/util"
)

const (
	MaxUsernameLength = 32 //用户名最大长度
	MaxPasswordLength = 32 //密码最大长度
	//MinPasswordLength = 8  //密码最小长度
)

// UserRegister 用户注册服务
func UserRegister(username string, password string) (common.UserIdTokenResp, error) {

	//0.数据准备
	var userResponse = common.UserIdTokenResp{}

	//1.合法性检验
	err := isLegal(username, password)
	if err != nil {
		return userResponse, err
	}

	//2.新建用户
	newUser, err := db.CreateUser(username, password)
	if err != nil {
		return userResponse, err
	}

	//3.颁发token
	token, err := middleware.CreateToken(newUser.ID, newUser.Name)
	if err != nil {
		log.Panicln("service-UserRegister: 创建用户token出错,", err)
		return userResponse, err
	}

	userResponse = common.UserIdTokenResp{
		UserId: newUser.ID,
		Token:  token,
	}
	return userResponse, nil
}

// UserInfo 用户信息获取服务
func UserInfo(rawId string, token string) (common.UserInfoResp, error) {
	// 数据准备
	var userInfoQueryResponse = common.UserInfoResp{}
	guestIDTemp, err := strconv.ParseUint(rawId, 10, 64)
	if err != nil {
		log.Panicln("service-UserInfo: 解析rawID时发生错误， ", err)
		return userInfoQueryResponse, err
	}
	guestID := uint(guestIDTemp)
	tokenStruct, ok := middleware.ParseToken(token)
	if err != nil {
		log.Panicln("解析token时发生错误")
	}
	hostID := tokenStruct.UserId

	// 获取用户信息, 先去查cache, cache查不到再查mysql
	user, ok := cache.GetUser(guestID)
	if !ok {
		user, err = db.GetUserInfo(guestID)
		if err != nil {
			return userInfoQueryResponse, err
		}
		cache.SetUser(&user)
	}

	//开始查询follow信息
	isFollowing, ok := cache.CheckFollowing(hostID, guestID)
	if !ok {
		isFollowing = db.CheckFollowing(hostID, guestID)
	}
	userInfoQueryResponse = util.PackUserInfo(user, isFollowing)
	return userInfoQueryResponse, nil
}

// UserLogin 用户登录服务
func UserLogin(username string, password string) (common.UserIdTokenResp, error) {

	// 数据准备
	var userResponse = common.UserIdTokenResp{}

	// 合法性检验
	err := isLegal(username, password)
	if err != nil {
		log.Println("用户名密码非法", err)
		return userResponse, err
	}

	// 避免缓存失效等操作, 用户登录等安全性较高信息不使用缓存
	// 检查用户是否存在
	user, flagExist := db.CheckUsername(username)
	if !flagExist {
		log.Println("Service-UserLogin: 登录失败: 用户 ", username, " 不存在.")
		return userResponse, errno.ErrorFullPossibility
	}

	// 检查密码是否正确
	if !checkPassword(user.Password, password) {
		log.Println("service-UserLogin: 登录失败：密码错误")
		return userResponse, errno.ErrorWrongPassword
	}

	if user.Model.ID == 0 {
		log.Println("账号或密码出错")
		return userResponse, errno.ErrorFullPossibility
	}

	// 颁发token
	token, err := middleware.CreateToken(user.Model.ID, user.Name)
	if err != nil {
		log.Println("service-UserLogin: 登录失败，创建用户token发生错误")
		return userResponse, err
	}
	//将用户信息写入cache
	cache.SetUser(&user)
	userResponse = common.UserIdTokenResp{
		UserId: user.Model.ID,
		Token:  token,
	}

	return userResponse, nil
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
		log.Panicln("用户名为空")
		return errno.ErrorNullUsername
	}
	if len(username) > MaxUsernameLength {
		log.Panicf("用户名过长，应小于%d位\n", MaxUsernameLength)
		return errno.ErrorUsernameExtend
	}

	//2.密码检验
	if password == "" {
		log.Panicln("密码为空")
		return errno.ErrorNullPassword
	}
	if len(password) > MaxPasswordLength {
		log.Panicf("密码过长，应小于%d位\n", MaxPasswordLength)
		return errno.ErrorPasswordLength
	}
	return nil
}
