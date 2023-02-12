package mapper

import (
	"tiktok/config"
	"tiktok/model"
	"time"
)

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
