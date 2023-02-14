package mapper

import (
	"errors"
	"log"
	"tiktok/config"
	"tiktok/model"
	"time"

	"gorm.io/gorm"
)

// IsFavorite 查询某用户是否点赞某视频
func IsFavorite(uid uint, vid uint) bool {
	if uid == 0 {
		//uid为0代表用户未登录, 默认返回false, 代表未关注
		return false
	}
	var total int64
	err := DBConn.Table("likes").Where("user_id = ? AND video_id = ? AND is_like = ?", uid, vid, true).Count(&total).Error
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
	create := DBConn.Table("videos").Create(&video)
	if create.Error != nil {
		//log.Println("视频信息插入数据库失败", create.Error)
		return create.Error
	}
	return nil

}
func FindVideosByUserID(userId uint) (resultVideos []model.Video, err error) {
	find := DBConn.Table("videos").Where("author_id=?", userId).Find(&resultVideos)
	if find.Error != nil {
		err = find.Error
	}
	return resultVideos, err
}

func FindVideosByLastTime(lastTime time.Time) (resultVideos []model.Video, err error) {
	find := DBConn.Table("videos").Where("created_at < ?", lastTime).Order("created_at desc").Limit(config.MaxFeedVideoCount).Find(&resultVideos)
	if find.Error != nil {
		err = find.Error
	}
	return resultVideos, err
}

// ExistVideo 检查videos表中是否存在vid对应的视频
func ExistVideo(vid uint) (video model.Video, flagExist bool) {
	err := DBConn.Table("videos").Where("ID = ?", vid).First(&video).Error
	return video, !errors.Is(err, gorm.ErrRecordNotFound)
}

// ReduceCommentCount 增加video记录中的评论计数
func AddCommentCount(vid uint) error {
	err := DBConn.Table("videos").Where("id = ?", vid).Update("comment_count", gorm.Expr("comment_count + 1")).Error
	if err != nil {
		log.Panicln("mapper-AddCommentCount: 增加视频的评论数量失败")
		return err
	}
	return nil
}

// ReduceCommentCount 减少video记录中的评论计数
func ReduceCommentCount(videoId uint) error {

	err := DBConn.Table("videos").Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count - 1")).Error
	if err != nil {
		return err
	}
	return nil
}

func GetVideoAuthor(videoId uint) (uint, error) {
	var video model.Video
	err := DBConn.Table("videos").Where("id = ?", videoId).Find(&video).Error
	if err != nil {
		return video.ID, err
	}
	return video.AuthorID, nil
}
