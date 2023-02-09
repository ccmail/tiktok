package controller

import (
	"net/http"
	"tiktok/pkg/common"
	middleware "tiktok/pkg/mw"
	"tiktok/service"

	"github.com/gin-gonic/gin"
)

// UserRegister 用户注册控制
func UserRegister(c *gin.Context) {
	//1.参数提取
	username, password := c.Query("username"), c.Query("password")

	//2.service层处理
	registerResponse, err := service.UserRegisterService(username, password)

	//3.返回响应
	if err != nil {
		c.JSON(http.StatusOK, common.UserRegisterResponse{
			BaseResponse: common.BaseResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
			UserIdTokenResponse: common.UserIdTokenResponse{},
		})
		return
	}
	c.JSON(http.StatusOK, common.UserRegisterResponse{
		BaseResponse: common.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "注册成功"},
		UserIdTokenResponse: registerResponse,
	})
}

// UserLogin 用户登录主函数
func UserLogin(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	userLoginResponse, err := service.UserLoginService(username, password)

	// 用户不存在返回对应的错误
	if err != nil {
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
	c.JSON(http.StatusOK, common.UserLoginResponse{
		BaseResponse: common.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "登录成功"},
		UserIdTokenResponse: userLoginResponse,
	})
}

// UserInfo 用户信息主函数
func UserInfo(c *gin.Context) {
	// 根据user_id查询
	rawId := c.Query("user_id")
	userInfoResponse, err := service.UserInfoService(rawId)

	// 根据token获得当前用户的userid
	token := c.Query("token")
	tokenStruct, _ := middleware.CheckToken(token)
	hostId := tokenStruct.UserId
	userInfoResponse.IsFollow = service.IsFollow(rawId, hostId)

	// 用户不存在返回对应的错误
	if err != nil {
		c.JSON(http.StatusOK, common.UserInfoResponse{
			BaseResponse: common.BaseResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
			UserList: common.UserInfoQueryResponse{},
		})
		return
	}

	// 用户存在，返回相应的id和token
	c.JSON(http.StatusOK, common.UserInfoResponse{
		BaseResponse: common.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "查询成功",
		},
		UserList: userInfoResponse,
	})

}
