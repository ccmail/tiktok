package util

import (
	"tiktok/model"
	"tiktok/pkg/common"
)

// PackVideoInfo 将数据库中的结构包装位响应体
func PackVideoInfo(videos model.Video, author common.UserInfoResp, isFavorite bool) common.VideoResp {
	returnVideo := common.VideoResp{
		ID:            videos.ID,
		Author:        author,
		PlayUrl:       videos.PlayUrl,
		CoverUrl:      videos.CoverUrl,
		FavoriteCount: videos.FavoriteCount,
		CommentCount:  videos.CommentCount,
		IsFavorite:    isFavorite,
		Title:         videos.Title,
	}
	return returnVideo
}

// PackUserInfo 包装user信息, 将其包装为response格式
func PackUserInfo(userInfo model.User, isFollow bool) common.UserInfoResp {
	userResp := common.UserInfoResp{
		UserID:        userInfo.ID,
		Username:      userInfo.Name,
		FollowCount:   userInfo.FollowCount,
		FollowerCount: userInfo.FollowerCount,
		IsFollow:      isFollow,
	}
	return userResp
}
