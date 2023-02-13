package mapper

import (
	"errors"
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
