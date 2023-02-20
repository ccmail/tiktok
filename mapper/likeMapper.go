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

func CreateLikeRecord(userID uint, videoID uint, isLike bool) error {
	likeRecord := model.Like{
		UserID:  userID,
		VideoID: videoID,
		IsLike:  isLike,
	}
	err := DBConn.Table("likes").Create(&likeRecord).Error
	if err != nil { //创建记录
		return err
	}
	if isLike {
		//DBConn.Table("videos").Where("id = ?", videoID).Update("like_count", gorm.Expr("like_count + 1"))
		DBConn.Table("videos").Where("id = ?", videoID).Update("favorite_count", gorm.Expr("favorite_count + 1"))
	}
	return nil
}

func UpdateLikeRecord(userID uint, videoID uint, isLike bool) {
	DBConn.Table("likes").Where("user_id = ? AND video_id = ?", userID, videoID).Update("is_like", isLike)
	if isLike {
		DBConn.Table("videos").Where("id = ?", videoID).Update("favorite_count", gorm.Expr("favorite_count + 1"))
	} else {
		DBConn.Table("videos").Where("id = ?", videoID).Update("favorite_count", gorm.Expr("favorite_count - 1"))
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

// CheckLikesNoHit 这里只查没有命中的部分, 之后直接写入到isFavorite中
func CheckLikesNoHit(hostID uint, isFavorite *[]bool, likeNoCache *map[uint][]int) (err error) {
	likes := make([]uint, 0, len(*likeNoCache)>>2)
	for k := range *likeNoCache {
		likes = append(likes, k)
	}

	var likeList []model.Like
	err = DBConn.Table("likes").Where("user_id = ? AND video_id IN ?", hostID, likes).Find(&likeList).Error
	if err != nil {
		log.Println("在mysql中查询点赞关系时失败")
		return err
	}

	for _, l := range likeList {
		for _, v := range (*likeNoCache)[l.VideoID] {
			(*isFavorite)[v] = l.IsLike
		}
	}
	return nil
}
