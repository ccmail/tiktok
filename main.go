package main

import (
	"tiktok/config"
	"tiktok/pkg/middleware"
)

func main() {
	// 连接数据库
	err := config.InitDBConnector()
	if err != nil {
		panic(err)
	}
	err = config.InitRedisConnector()
	if err != nil {
		panic(err)
	}
	config.InitLog()
	// 连接OSS服务
	err = middleware.InitOSS()
	if err != nil {
		panic(err)
	}
	//注册路由
	r := config.InitRouter()
	//启动端口为12345的项目
	errRun := r.Run(":12345")
	if errRun != nil {
		return
	}
}

// func main() {
// 	strLastTime := ""
// 	lastTime, err := strconv.ParseInt(strLastTime, 10, 32)
// 	if err != nil {
// 		lastTime = 0
// 	}
// 	fmt.Println(lastTime)
// }
