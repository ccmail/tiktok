package service

import (
	"log"
	"tiktok/mapper"
	"tiktok/model"
	"tiktok/pkg/common"
)

// LikeService
func LikeService(userID uint, videoID uint, actionType uint) error {
	//如果没有记录-Create，如果有了记录-修改IsLike
	likeExist, flagExist := mapper.ExistLikeRecord(userID, videoID)
	if !flagExist { // 不存在记录
		err := mapper.CreateLikeRecord(userID, videoID, true)
		if err != nil {
			log.Panicln("service-LikeService: 创建like记录失败，", err)
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
		if !likeExist.IsLike { //IsLike为false，则video的like_count加1
			mapper.UpdateLikeRecord(userID, videoID, true)
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
	likeRecord, flagExist := mapper.ExistLikeRecord(userID, videoID)
	if !flagExist { // 不存在记录
		err := mapper.CreateLikeRecord(userID, videoID, false)
		if err != nil {
			log.Panicln("service-LikeService: 创建like记录失败，", err)
			return err
		}
	} else { // 存在记录
		if likeRecord.IsLike { //IsLike为ture，则video的like_count减1
			mapper.UpdateLikeRecord(userID, videoID, false)

		}
		//IsLike为false-video的like_count不变
	}

	return nil
}

// GetLikeList 获取点赞列表
func LikeListService(userID uint) ([]model.Video, error) {

	//查询当前id用户的所有点赞视频

	videoList, err := mapper.GetLikeList(userID)

	if err != nil {
		log.Panicln("service-LikeListServic: 获取喜欢列表失败，", err)
	}

	return videoList, nil
}

// 为要返回的视频列表填充信息
func FillInfo(videoList []model.Video, userIdHost uint) []common.FavoriteVideo {
	returnList := make([]common.FavoriteVideo, 0)
	for _, m := range videoList {
		var author = common.UserInfoQueryResponse{}
		var getAuthor = model.User{}
		getAuthor, err := mapper.FindUserInfo(m.AuthorID)
		if err != nil {
			log.Println("未找到作者: ", m.AuthorID)
			continue
		}
		// 是否已关注
		isfollowing := mapper.CheckFollowing(userIdHost, m.AuthorID)
		// 是否已点赞
		isfavorite := mapper.IsFavorite(userIdHost, m.ID)

		// 作者信息
		author.UserID = getAuthor.ID
		author.Username = getAuthor.Name
		author.FollowCount = getAuthor.FollowCount
		author.FollowerCount = getAuthor.FollowerCount
		author.IsFollow = isfollowing

		video := common.FavoriteVideo{
			ID:            m.ID,
			Author:        author,
			PlayUrl:       m.PlayUrl,
			CoverUrl:      m.CoverUrl,
			FavoriteCount: m.LikeCount,
			CommentCount:  m.CommentCount,
			IsFavorite:    isfavorite,
			Title:         m.Title,
		}
		returnList = append(returnList, video)
	}
	return returnList
}
