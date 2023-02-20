package service

import (
	"log"
	"tiktok/mapper/gorm"
	"tiktok/model"
	"tiktok/pkg/common"
	"tiktok/pkg/errno"
	"tiktok/pkg/util"
)

// LikeService 这里逻辑可能有点问题, actionType形参没有使用
func LikeService(userID uint, videoID uint, actionType uint) error {
	// 首先要保证视频存在
	_, videoExist := gorm.ExistVideo(videoID)
	if !videoExist {
		log.Println("service-LikeService: 点赞失败，未找到对应视频")
		return errno.ErrorNullVideo
	}
	//如果没有记录-Create，如果有了记录-修改IsLike
	likeRecord, likeExist := gorm.ExistLikeRecord(userID, videoID)
	if !likeExist { // 不存在记录
		err := gorm.CreateLikeRecord(userID, videoID, true)
		if err != nil {
			log.Println("service-LikeService: 创建like记录失败，", err)
			return err
		}
		// //userId的like_count增加
		// if err := mapper.AddLikeCount(userID); err != nil {
		// 	return err
		// }
		// //videoId对应的userId的total_like增加
		// GuestId, err := GetVideoAuthor(videoID)
		// if err != nil {
		// 	return err
		// }
		// if err := mapper.AddTotalLiked(GuestId); err != nil {
		// 	return err
		// }
	} else { // 存在记录
		if !likeRecord.IsLike { //IsLike为false，则video的like_count加1
			gorm.UpdateLikeRecord(userID, videoID, true)
			// //userId的like_count增加
			// if err := mapper.AddLikeCount(userID); err != nil {
			// 	return err
			// }
			// //videoId对应的userId的total_like增加
			// GuestId, err := GetVideoAuthor(videoID)
			// if err != nil {
			// 	return err
			// }
			// if err := mapper.AddTotalLiked(GuestId); err != nil {
			// 	return err
			// }
		}
		// IsLike为true则无需处理
		return nil
	}

	return nil
}

func CancelLikeService(userID uint, videoID uint) error {
	// 首先要保证视频存在
	_, videoExist := gorm.ExistVideo(videoID)
	if !videoExist {
		log.Panicln("service-LikeService: 点赞失败，未找到对应视频")
		return errno.ErrorNullVideo
	}
	likeRecord, likeExist := gorm.ExistLikeRecord(userID, videoID)
	if !likeExist { // 不存在记录
		err := gorm.CreateLikeRecord(userID, videoID, false)
		if err != nil {
			log.Panicln("service-LikeService: 创建like记录失败，", err)
			return err
		}
	} else { // 存在记录
		if likeRecord.IsLike { //IsLike为ture，则video的like_count减1
			gorm.UpdateLikeRecord(userID, videoID, false)

		}
		//IsLike为false-video的like_count不变
	}

	return nil
}

// LikeListService  获取点赞列表
func LikeListService(userID uint) ([]model.Video, error) {

	//查询当前id用户的所有点赞视频

	videoList, err := gorm.GetLikeList(userID)

	if err != nil {
		log.Panicln("service-LikeListService: 获取喜欢列表失败，", err)
	}

	return videoList, nil
}

// FillInfo 为要返回的视频列表填充信息
func FillInfo(videoList []model.Video, userIdHost uint) []common.VideoResp {
	returnList := make([]common.VideoResp, 0)
	for _, m := range videoList {
		var author = common.UserInfoResp{}
		var getAuthor = model.User{}
		getAuthor, err := gorm.FindUserInfo(m.AuthorID)
		if err != nil {
			log.Println("未找到作者: ", m.AuthorID)
			continue
		}

		// 作者信息
		author.UserID = getAuthor.ID
		author.Username = getAuthor.Name
		author.FollowCount = getAuthor.FollowCount
		author.FollowerCount = getAuthor.FollowerCount
		author.IsFollow = gorm.CheckFollowing(userIdHost, m.AuthorID)

		video := util.PackVideoInfo(m, author, gorm.IsFavorite(userIdHost, m.ID))
		/*		video := common.VideoResp{
				ID:            m.ID,
				Author:        author,
				PlayUrl:       m.PlayUrl,
				CoverUrl:      m.CoverUrl,
				FavoriteCount: m.FavoriteCount,
				CommentCount:  m.CommentCount,
				IsFavorite:    mapper.IsFavorite(userIdHost, m.ID),
				Title:         m.Title,
			}*/
		returnList = append(returnList, video)
	}
	return returnList
}
