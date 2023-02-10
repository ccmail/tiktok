# service
service层用于实现业务逻辑


## zxl, 2023.2.9
新增 userService.go，处理用户注册/登录/获取用户信息接口的相关逻辑

## zxl, 2023.2.10
新增 videoService.go，处理投稿、视频流、发布列表接口的相关逻辑\
新增 followerService.go，处理用户关注相关接口，目前只实现了一个检查两个用户是否为关注关系的函数。