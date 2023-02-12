本文的用于记录程季晟的一些改动和疑问

# 待办

1. `videoService`中上传云服务器等内容的拆分, 目前推荐放到`util`工具包下比较合适, **暂定**

2. ~~video方面的feed()内容~~
3. ~~get请求获取不到token,不清楚原因~~

# 改动

## 2023.2.10

1. 更改了一些注释
2. 在routers根路由处添加"douyin"字段(第12行)

## 2023.2.11

1. `/pkg/mw/`更名为`pkg/middleware`
2. ~~`/pkg/middleware/auth`中, token加密方式由HS256更改为RS256(非对称加密)~~
3. `/pkg/middleware/auth`中加密密钥`Key`删除, 在`config/config`中添加了对应的`Key`常量, 注意, 常量为字符串,
   使用时需要转化为[]byte
4. `/pkg/middleware/auth`中, 过期时间`24`删除, 在`config/config`添加了对应的`TokenLiveTime`常量, 注意, 该常量为int类型,
   使用时需要转化为`time.Duration`
5. `/pkg/middleware/`中, `CheckToken()`更名为`ParseToken()`, 已经引用`CheckToken()`的部分已经使用ide安全重构
6. `/service/videoService`中, 将`ExampleReadFrameAsJpeg()`的异常拦截, 交由使用该服务的部分判断
7. **将查找用户信息的逻辑拆分, 直接查库的移动到了`userMapper.go`中**
8. **publish功能中对应的Create插入语句, 移动到了对应的mapper中**
9. **对7,8两条对应的数据库操作进行了封装, 并将error作为了返回项, 用以在service层拦截错误**
10. `pkg/common/baseResponse.go:57`中的`ReturnVideo`更正了json映射, 使其与api一一对应

## 2023.2.12

1. 重构了common包中的结构, 增加了可复用性

2. video的相关操作, 如`Feed, Publish, PublishList`基本完成

3. 将数据库字段对应的数据结构, 包装成返回给客户端response这一过程, 抽象成了方法

4. 将`ffmpeg`, `oos`等操作移动到`middleWare`中

5. **取消了`PublishList`的鉴权操作**, 业务场景如下:

   > 1. 用户未登录, 希望访问某视频作者的个人主页
        >
   2. 用户已经登录, 希望查看自己的主页

   因为`publishList`请求时具有`user_id`, 因此暂时, 两个操作可以复用一个方法

# 新增

## 2023.2.11

1. `config/config.go`文件, 用于定义一些常量
2. `pkg/middleware/log.go`文件, 这里思路没有定好, 打算在这里定义log设置, 不知道可不可行
3. `controller/videoController_cjs.go`实现了`publish`和`publishList`
4. `service/videoService_cjs.go`和`service/videoServiceImpl.go`, 定义了`VideoServiceImpl`接口, 并在`videoServiceImpl`
   中实现对应接口, 如果其他controller层需要修改, 请参考 `videoController_cjs.go`和`videoServiceImpl.go`等文件的写法
5. 若干test文件, **建议test文件不要删除, 推测是加分项**
6. `pkg/middleware/auth.go:66`新增了`parseToken`返回error的写法, **推荐**各个方法尽量返回err, 用以拦截错误信息,
   以保证程序的稳定性

## 2023.2.12

1. `util/removeIllegalChar.go:6`, 因为上传oss过程中, 存储名称存在" "等字符会报错, 因此使用其去除文件名中的非法字符,
   用以存储到oss中, 非法字符的定义详见`config/config.go:10`
2. 将从数据库中查到的字段, 封装为响应体的这一过程抽象为方法, 实现复用. 例如`publishList`和`Feed`中, 部分方法可以复用

# 疑问

1. 是否需要严格按照MVC模式进行分层, 即Controller用于拦截, 响应请求, Service用于实现逻辑(
   通过定义接口和实现接口方法完成逻辑)

2. 如果按照1执行,现在是先实现功能后期再改, 还是初步实现时就

    