package controller

import (
	"log"
	"net/http"
	"tiktok/pkg/common"
	"tiktok/pkg/middleware"

	"tiktok/service"

	"github.com/gin-gonic/gin"
)

// UserRegister 用户注册控制层
func UserRegister(c *gin.Context) {
	//1.参数提取
	username, password := c.Query("username"), c.Query("password")

	//2.service层处理
	registerResponse, err := service.UserRegisterService(username, password)

	//3.返回响应
	if err != nil {
		log.Println("controller-UserRegister: 注册失败: 用户名: ", username, ", 密码: ", password)
		c.JSON(http.StatusOK, common.UserRegisterResponse{
			BaseResponse: common.BaseResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
			UserIdTokenResponse: common.UserIdTokenResponse{},
		})
		return
	}
	log.Println("用户", username, "注册成功。")
	c.JSON(http.StatusOK, common.UserRegisterResponse{
		BaseResponse: common.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "注册成功"},
		UserIdTokenResponse: registerResponse,
	})
}

// UserLogin 用户登录控制层
func UserLogin(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	userLoginResponse, err := service.UserLoginService(username, password)

	// 用户不存在返回对应的错误
	if err != nil {
		log.Println("controller-UserInfo: 用户登录失败,", err)
		c.JSON(http.StatusOK, common.UserLoginResponse{
			BaseResponse: common.BaseResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
			UserIdTokenResponse: common.UserIdTokenResponse{},
		})
		return
	}

	// 用户存在，返回相应的id和token
	log.Println("用户", username, "登录成功。")
	c.JSON(http.StatusOK, common.UserLoginResponse{
		BaseResponse: common.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "登陆成功"},
		UserIdTokenResponse: userLoginResponse,
	})
}

// UserInfo 用户信息控制层
func UserInfo(c *gin.Context) {
	// 根据user_id查询
	rawId := c.Query("user_id")
	userInfoResponse, err := service.UserInfoService(rawId)

	// 用户不存在返回对应的错误
	if err != nil {
		log.Println("controller-UserInfo: 查找用户信息失败", err)
		c.JSON(http.StatusOK, common.UserInfoResponse{
			BaseResponse: common.BaseResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
			UserList: common.UserInfoQueryResponse{},
		})
		return
	}

	// 根据token获得当前用户的userid
	token := c.Query("token")
	tokenStruct, _ := middleware.CheckToken(token)
	hostId := tokenStruct.UserId
	userInfoResponse.IsFollow = service.IsFollow(rawId, hostId)

	// 用户存在，返回相应的id和token
	log.Println("获取用户id", rawId, "的信息成功。")
	c.JSON(http.StatusOK, common.UserInfoResponse{
		BaseResponse: common.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "获取用户信息成功",
		},
		UserList: userInfoResponse,
	})

}
