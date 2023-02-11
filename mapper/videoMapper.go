package mapper

import (
	"tiktok/model"
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
func FindVideoList(userId uint) ([]model.Video, error) {
	var videoList []model.Video
	find := DBConn.Table("videos").Where("author_id=?", userId).Find(&videoList)
	if find.Error != nil {
		return videoList, find.Error
	}
	return videoList, nil
}
