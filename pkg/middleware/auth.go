package middleware

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"tiktok/config"
	"tiktok/pkg/common"
	"time"
)

//var Key = []byte("byte dance 11111 return")

// MyClaims
// jwt生成/解析token的格式, 用来传递/接受token中的信息
type MyClaims struct {
	UserId   uint   `json:"user_id"`
	UserName string `json:"username"`
	jwt.StandardClaims
}

// CreateToken
// 生成token
func CreateToken(userId uint, userName string) (string, error) {
	expireTime := time.Now().Add(time.Duration(config.TokenLiveTime) * time.Hour) //过期时间
	nowTime := time.Now()                                                         //当前时间
	claims := MyClaims{
		UserId:   userId,
		UserName: userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间戳
			IssuedAt:  nowTime.Unix(),    //当前时间戳
			Issuer:    "henrik",          //颁发者签名
			Subject:   "userToken",       //签名主题
		},
	}
	//更改了加密方式, 更改为非对称加密RS256
	tokenStruct := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用密钥加密token, 传入要求[]byte切片, 记得类型转换
	return tokenStruct.SignedString([]byte(config.Key))
}

// ParseToken
// 解析Token, 这里可能需要把名字更改为ParseToken
// 返回值为颁发的token中所携带的信息,以及是否查到token
func ParseToken(token string) (*MyClaims, bool) {
	tokenObj, err := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Key), nil
	})
	if err != nil {
		log.Panicln("token无效, 请检查token", err)
		return nil, false
	}
	if key, _ := tokenObj.Claims.(*MyClaims); tokenObj.Valid {
		return key, true
	} else {
		return nil, false
	}
}

// ParseTokenCJS
// 解析Token, 这里可能需要把名字更改为ParseToken
// 返回值为颁发的token中所携带的信息,以及是否查到token
func ParseTokenCJS(token string) (*MyClaims, error) {
	tokenObj, err := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Key), nil
	})
	if err != nil {
		log.Panicln("token无效, 请检查token", err)
		return nil, err
	}
	if key, _ := tokenObj.Claims.(*MyClaims); tokenObj.Valid {
		return key, nil
	} else {
		return nil, errors.New("token解析失败")
	}
}

// JwtMiddleware jwt 中间件
func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Query("token")
		if tokenStr == "" {
			tokenStr = c.PostForm("token")
		}
		// 用户不存在
		if tokenStr == "" {
			c.JSON(http.StatusOK, common.BaseResponse{
				StatusCode: 401,
				StatusMsg:  "Target user not found",
			},
			)
			c.Abort() // 阻止执行
			return
		}
		// 验证 token
		tokenStruck, ok := ParseToken(tokenStr)
		if !ok {
			c.JSON(http.StatusOK, common.BaseResponse{
				StatusCode: 403,
				StatusMsg:  "Wrong token",
			})
			c.Abort() // 阻止执行
			return
		}
		// token 超时
		if time.Now().Unix() > tokenStruck.ExpiresAt {
			c.JSON(http.StatusOK, common.BaseResponse{
				StatusCode: 402,
				StatusMsg:  "Expired token",
			})
			c.Abort() // 阻止执行
			return
		}
		c.Set("username", tokenStruck.UserName)
		c.Set("user_id", tokenStruck.UserId)

		c.Next()
	}
}
