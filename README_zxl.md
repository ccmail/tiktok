# 接口开发进度

## 1 基础接口

- [x] 视频流接口：GET /douyin/feed/
- [x] 用户注册：POST /douyin/user/register/
- [x] 用户登录：POST /douyin/user/login/
- [x] 用户信息：GET /douyin/user/
- [x] 投稿接口 POST /douyin/publish/action/
- [x] 发布列表 GET /douyin/publish/list/

## 2 互动接口

- [ ] 赞操作：POST /douyin/favorite/action/
- [ ] 喜欢列表：GET /douyin/favorite/list/
- [ ] 评论操作：POST /douyin/comment/action/
- [ ] 评论列表：GET /douyin/comment/list/

## 3 社交接口

- [ ] 关注操作：POST /douyin/relation/action/
- [ ] 关注列表：GET /douyin/relation/follow/list/
- [ ] 粉丝列表：GET /douyin/relation/follower/list/
- [ ] 好友列表：GET /douyin/relation/friend/list/
- [ ] 发送消息：POST /douyin/message/action/
- [ ] 聊天记录：GET /douyin/message/chat/

# zxl, 2023.2.9
实现了用户注册、用户登录、用户查询接口

## DAO 部分的改动
参见 ReadMeModel.md

## 项目根目录的改动

- go.mod，module 名改为 tiktok （一般都用小写）
- 新增 config 目录，包含 db.yaml 和 oss.yaml 分别作为数据库和OSS的配置文件。为防隐私泄露，这两个文件加入了.gitignore中。
- 新增 pkg/errno 目录，用于定义错误类型
- 新增 pkg/mw 目录，用于存放中间件（例如 jwt 鉴权）
- 新增 pkg/common 目录，用于定义所有类型的 response

# zxl, 2023.2.10

实现了投稿、发布列表、视频流接口

# zxl 2023.2.12

## 用户注册、登录、用户信息接口重构完成

- 将userService中的dao部分相关函数移到了mapper/userMapper.go中，包括
    - GetUser函数改名为cjs实现的FindUserInfo
    - CreateUser函数（所需的encrypt也做了移动）
    - ExistUsername函数
- 将在DB中检查关注关系的IsFollowing函数改名为CheckFollowing函数（因为可能与IsFollow函数产生歧义）并移到mapper/followerMapper.go中
- 错误处理相关代码添加日志信息
- 将原先的英文错误信息、Response中英文StatusMsg改为中文

> video部分的未做改动，待后期merge

## 实现点赞、喜欢列表接口

- 将viedel实体的FavoriteCount字段改名为LikeCount
- 目录结构对应cjs分支的做相应改动。