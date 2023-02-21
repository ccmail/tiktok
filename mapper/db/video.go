package db

import (
	"errors"
	"log"
	"tiktok/config"
	"tiktok/mapper"
	"tiktok/model"
	"time"

	"gorm.io/gorm"
)

type Option func(db *gorm.DB)

func GetVideoInfo(opts ...Option) (v []*model.Video, err error) {
	//优先去查redis
	for i := range opts {
		opts[i](mapper.DBConn)
	}
	if err = mapper.DBConn.Error; err != nil {
		log.Panicln("查询失败")
		return v, err
	}
	mapper.DBConn.Find(&v)
	return v, nil
}

func VideoIDList(userID uint) Option {
	return func(db *gorm.DB) {
		db.Where("author_id=?", userID)
	}
}

// CheckIsFavorite 查询某用户是否点赞某视频
func CheckIsFavorite(uid uint, vid uint) bool {
	if uid == 0 {
		//uid为0代表用户未登录, 默认返回false, 代表未关注
		return false
	}
	var total int64
	err := mapper.DBConn.Table("likes").Where("user_id = ? AND video_id = ? AND is_like = ?", uid, vid, true).Count(&total).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	if total == 0 {
		return false
	}
	return true
}

// CreateVideo 添加一条视频信息
func CreateVideo(video *model.Video) error {
	create := mapper.DBConn.Table("videos").Create(&video)
	if create.Error != nil {
		//log.Println("视频信息插入数据库失败", create.Error)
		return create.Error
	}
	//根据作者id添加视频到cache
	//SetVideoCache(*video, video.AuthorID)
	return nil

}
func GetVideosByUserID(userId uint) (resultVideos []model.Video, err error) {
	//resultVideos = GetVideoCache(userId)
	//if len(resultVideos) <= 0 {
	//	log.Println("cache中没有该用户的相关视频")
	//}
	find := mapper.DBConn.Table("videos").Where("author_id=?", userId).Find(&resultVideos)
	if find.Error != nil {
		err = find.Error
	}
	return resultVideos, err
}

// FindVideosByLastTime 这里需要为所有的视频设置一个list, list长度为30
func FindVideosByLastTime(lastTime time.Time) (resultVideos []model.Video, err error) {

	find := mapper.DBConn.Table("videos").Where("created_at < ?", lastTime).Order("created_at desc").Limit(config.MaxFeedVideoCount).Find(&resultVideos)
	if find.Error != nil {
		err = find.Error
	}
	return resultVideos, err
}

// CheckVideo 检查videos表中是否存在vid对应的视频
func CheckVideo(vid uint) (video model.Video, flagExist bool) {
	err := mapper.DBConn.Table("videos").Where("id = ?", vid).First(&video).Error
	return video, !errors.Is(err, gorm.ErrRecordNotFound)
}

// AddCommentCount  增加video记录中的评论计数
func AddCommentCount(vid uint) error {
	err := mapper.DBConn.Table("videos").Where("id = ?", vid).Update("comment_count", gorm.Expr("comment_count + 1")).Error
	if err != nil {
		log.Panicln("mapper-AddCommentCount: 增加视频的评论数量失败")
		return err
	}
	return nil
}

// ReduceCommentCount 减少video记录中的评论计数
func ReduceCommentCount(videoId uint) error {

	err := mapper.DBConn.Table("videos").Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count - 1")).Error
	if err != nil {
		return err
	}
	return nil
}

func GetVideo(videoId uint) (model.Video, error) {
	var video model.Video
	err := mapper.DBConn.Table("videos").Where("id = ?", videoId).Find(&video).Error
	if err != nil {
		return video, err
	}
	return video, nil
}
