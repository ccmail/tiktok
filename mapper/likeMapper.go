package mapper

import (
	"errors"
	"log"
	"tiktok/model"

	"gorm.io/gorm"
)

// ExistLikeRecord true-存在记录，false-不存在记录
func ExistLikeRecord(userId uint, videoId uint) (likeRecord model.Like, flagExist bool) {
	err := DBConn.Table("likes").Where("user_id = ? AND video_id = ?", userId, videoId).First(&likeRecord).Error
	return likeRecord, !errors.Is(err, gorm.ErrRecordNotFound)
}

func CreateLikeRecord(userID uint, videoID uint, islike bool) error {
	likeRecord := model.Like{
		UserID:  userID,
		VideoID: videoID,
		IsLike:  islike,
	}
	err := DBConn.Table("likes").Create(&likeRecord).Error
	if err != nil { //创建记录
		return err
	}
	if islike {
		DBConn.Table("videos").Where("id = ?", videoID).Update("like_count", gorm.Expr("like_count + 1"))
	}
	return nil
}

func UpdateLikeRecord(userID uint, videoID uint, islike bool) {
	DBConn.Table("likes").Where("user_id = ? AND video_id = ?", userID, videoID).Update("is_like", islike)
	if islike {
		DBConn.Table("videos").Where("id = ?", videoID).Update("like_count", gorm.Expr("like_count + 1"))
	} else {
		DBConn.Table("videos").Where("id = ?", videoID).Update("like_count", gorm.Expr("like_count - 1"))
	}
}

func GetLikeList(userID uint) (videoList []model.Video, err error) {
	var likeList []model.Like
	videoList = make([]model.Video, 0)
	err = DBConn.Table("likes").Where("user_id=? AND is_like=?", userID, true).Find(&likeList).Error
	if err != nil { // 找不到记录
		log.Println("mapper-GetLikeList: 未找到喜欢的视频")
	}
	for _, m := range likeList {
		var video = model.Video{}
		if err := DBConn.Table("videos").Where("id=?", m.ID).Find(&video).Error; err != nil {
			return nil, err
		}
		videoList = append(videoList, video)
	}
	return videoList, nil
}

// func AddLikeCount(HostId uint) error {
// 	if err := DBConn.Model(&model.User{}).
// 		Where("id=?", HostId).
// 		Update("favorite_count", gorm.Expr("favorite_count+?", 1)).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

// func AddTotalLiked(HostId uint) error {
// 	if err := DBConn.Model(&model.User{}).
// 		Where("id=?", HostId).
// 		Update("total_favorited", gorm.Expr("total_favorited+?", 1)).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

// func ReduceLikeCount(HostId uint) error {
// 	if err := DBConn.Model(&model.User{}).
// 		Where("id=?", HostId).
// 		Update("favorite_count", gorm.Expr("favorite_count-?", 1)).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

// func ReduceTotalLiked(HostId uint) error {
// 	if err := DBConn.Model(&model.User{}).
// 		Where("id=?", HostId).
// 		Update("total_favorited", gorm.Expr("total_favorited-?", 1)).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }
