package util

import (
	"log"
	"tiktok/pkg/middleware"
)

// GetHostIDFromToken 从token返回id, 只有不强制需要使用token鉴权的方法才能用!!
func GetHostIDFromToken(hostToken string) uint {
	var hostID uint
	if hostToken != "" {
		hostInfo, err := middleware.ParseTokenCJS(hostToken)
		if err != nil {
			//就算token解析失败也不应该拒绝访问publishList
			log.Println("请求作品列表的用户的token解析失败", err)
		} else {
			hostID = hostInfo.UserId
		}
	}
	//log.Println(hostID)
	return hostID
}
