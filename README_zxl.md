# 接口开发进度

## 1 基础接口

- [ ] 视频流接口：GET /douyin/feed/

- [x] 用户注册：POST /douyin/user/register/
- [x] 用户登录：POST /douyin/user/login/
- [x] 用户信息：GET /douyin/user/
- [ ] 投稿接口 POST /douyin/publish/action/
- [ ] 发布列表 GET /douyin/publish/list/

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


## DAO 部分的改动
参见 ReadMeModel.md

## 项目根目录的改动
- go.mod，module 名改为 tiktok （一般都用小写）
- 新增 config 目录，包含 db.yaml 和 oss.yaml 分别作为数据库和OSS的配置文件。为防隐私泄露，这两个文件加入了.gitignore中。
- 新增 pkg/errno 目录，用于定义错误类型
- 新增 pkg/mw 目录，用于存放中间件（例如 jwt 鉴权）