package main

import (
	"tiktok/controller"
	"tiktok/mapper"
	"tiktok/pkg/middleware"
	"tiktok/pkg/util"
)

func main() {
	// 连接数据库
	err := mapper.InitDBConnector()
	if err != nil {
		panic(err)
	}
	err = mapper.InitRedisConnector()
	if err != nil {
		panic(err)
	}
	util.InitLog()
	// 连接OSS服务
	err = middleware.InitOSS()
	if err != nil {
		panic(err)
	}
	//注册路由
	r := controller.InitRouter()
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
