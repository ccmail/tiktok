# model

model层用于存放数据库字段对应的数据结构


## zxl, 2023.2.9
将实体定义从 mapper 包中移到 model 包中


- User 实体
  - 添加 gorm.Model 字段
  - UserName 字段改名为 Name 字段（与接口的json字段保持一致）
  - 删去 UserHeadPicURL （头像）字段（接口没有要求这个）
  - 删去 account 字段（注册/登录等接口都只需要username）
- Video 实体
  - VideoURL 字段改名为 PlayUrl 字段（与接口的json字段保持一致）
  - CoverPicURL字段改为 CoverUrl 字段（与接口的json字段保持一致）
  - 删去 CreatedAt, UpdatedAt, DeletedAt 字段（已包含在 gorm.Model字段中）
  - 新增 FavoriteCount （视频的点赞总数）、 comment_count （视频的评论总数）字段（接口要求）
- Follower 实体
- 添加 gorm.Model 字段
- 删去 FollowerHeadPicURL 字段

其余实体，因目前实现的接口暂未涉及到，未作改动