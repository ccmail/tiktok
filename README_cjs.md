本文的用于记录程季晟的一些改动和疑问

# 改动

## 2023.2.10

1. 更改了一些注释
2. 在routers根路由处添加"douyin"字段(第12行)

## 2023.2.11

1. `/pkg/mw/`更名为`pkg/middleware`
2. `/pkg/middleware/auth`中, token加密方式由HS256更改为RS256(非对称加密)
3. `/pkg/middleware/auth`中加密密钥`Key`删除, 在`config/config`中添加了对应的`Key`常量, 注意, 常量为字符串,
   使用时需要转化为[]byte
4. `/pkg/middleware/auth`中, 过期时间`24`删除, 在`config/config`添加了对应的`TokenLiveTime`常量, 注意, 该常量为int类型,
   使用时需要转化为`time.Duration`
5. `/pkg/middleware/`中, `CheckToken()`更名为`ParseToken()`, 已经引用`CheckToken()`的部分已经使用ide安全重构
6. `/service/videoService`中, 将`ExampleReadFrameAsJpeg()`的异常拦截, 交由使用该服务的部分判断

# 新增

## 2023.2.11

1. `config/config.go`文件, 用于定义一些常量
2. `pkg/middleware/log.go`文件, 这里思路没有定好, 打算在这里定义log设置, 不知道可不可行

# 疑问

1. 是否需要严格按照MVC模式进行分层, 即Controller用于拦截, 响应请求, Service用于实现逻辑(
   通过定义接口和实现接口方法完成逻辑)
2. 如果按照1执行,现在是先实现功能后期再改, 还是初步实现时就