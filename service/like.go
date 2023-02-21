package service

import (
	"context"
	"encoding/json"
	"log"
	"tiktok/config"
	"tiktok/mapper/cache"
	"tiktok/mapper/db"
	"tiktok/model"
	"tiktok/pkg/common"
	"tiktok/pkg/errno"
	"tiktok/pkg/mq"
	"tiktok/pkg/util"
)

// Like 这里逻辑可能有点问题, actionType形参没有使用
func Like(userID uint, videoID uint, actionType uint) error {
	isLike := actionType == 1
	// 首先要保证视频存在
	setVideo, videoExist := cache.GetVideo(videoID)
	if !videoExist {
		setVideo, videoExist = db.CheckVideo(videoID)
		if !videoExist {
			log.Println("service-Like: 点赞失败，未找到对应视频")
			return errno.ErrorNullVideo
		}
	}
	//更新cache
	cache.SetVideo(&setVideo)

	//先去cache中查询点赞记录
	likeRecord, likeExist := cache.CheckFavorite(userID, &setVideo)
	if !likeExist {
		likeRecord, likeExist = db.CheckLikeRecord(userID, videoID)
		if !likeExist { // 不存在记录
			//使用goroutine启动点赞操作
			marshal, err := json.Marshal(mq.LikeStruct{UserID: userID, VideoID: videoID, IsLike: isLike})
			if err != nil {
				log.Println("序列化失败")
			}
			go func() {
				err := mq.Q.SendMsg(context.Background(), &mq.Msg{
					ID:       "1",
					Topic:    config.MQLikeAdd,
					Body:     marshal,
					Group:    "g1",
					Consumer: "c1",
				})
				if err != nil {
					log.Panicln("发送消息失败!")
				}
			}()
		}
	}
	//点赞状态不同时更改
	if likeExist && likeRecord != isLike {
		//go mq.PubUpdateLike(userID, videoID, isLike)
		marshal, err := json.Marshal(mq.LikeStruct{UserID: userID, VideoID: videoID, IsLike: isLike})
		if err != nil {
			log.Println("序列化失败")
		}
		go func() {
			err := mq.Q.SendMsg(context.Background(), &mq.Msg{
				ID:       "1",
				Topic:    config.MQLikeUpdate,
				Body:     marshal,
				Group:    "g1",
				Consumer: "c1",
			})
			if err != nil {
				log.Panicln("发送消息失败!")
			}
		}()
		//go db.UpdateLikeRecord(userID, videoID, isLike)
	}

	//最后写入cache
	cache.SetFavorite(userID, videoID, isLike)

	return nil
}

// LikeList  获取点赞列表
func LikeList(userID uint) ([]model.Video, error) {

	//查询当前id用户的所有点赞视频

	videoList, err := db.GetLikeList(userID)

	if err != nil {
		log.Panicln("service-LikeList: 获取喜欢列表失败，", err)
	}

	return videoList, nil
}

// FillInfo 为要返回的视频列表填充信息
func FillInfo(videoList []model.Video, userIdHost uint) []common.VideoResp {
	returnList := make([]common.VideoResp, 0)
	for _, m := range videoList {
		var author = common.UserInfoResp{}
		var getAuthor = model.User{}
		getAuthor, err := db.GetUserInfo(m.AuthorID)
		if err != nil {
			log.Println("未找到作者: ", m.AuthorID)
			continue
		}

		// 作者信息
		author.UserID = getAuthor.ID
		author.Username = getAuthor.Name
		author.FollowCount = getAuthor.FollowCount
		author.FollowerCount = getAuthor.FollowerCount
		author.IsFollow = db.CheckFollowing(userIdHost, m.AuthorID)

		video := util.PackVideoInfo(m, author, db.CheckIsFavorite(userIdHost, m.ID))
		/*		video := common.VideoResp{
				ID:            m.ID,
				Author:        author,
				PlayUrl:       m.PlayUrl,
				CoverUrl:      m.CoverUrl,
				FavoriteCount: m.FavoriteCount,
				CommentCount:  m.CommentCount,
				CheckIsFavorite:    mapper.CheckIsFavorite(userIdHost, m.ID),
				Title:         m.Title,
			}*/
		returnList = append(returnList, video)
	}
	return returnList
}
