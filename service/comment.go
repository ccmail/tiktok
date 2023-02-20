package service

import (
	"log"
	"tiktok/mapper/cache"
	"tiktok/mapper/db"
	"tiktok/model"
	"tiktok/pkg/common"
	"tiktok/pkg/util"
)

// PostComment 发布评论
func PostComment(userId uint, text string, videoId uint) (model.Comment, error) {
	newComment := model.Comment{
		VideoID:     videoId,
		UserID:      userId,
		CommentText: text,
		Valid:       true,
	}

	err := db.NewCommentTx(newComment)
	if err != nil {
		return model.Comment{}, err
	}
	//发布之后, 直接插入缓存
	cache.SetComment(&newComment)
	return newComment, nil
}

func GetCommenter(userId uint) (commenter model.User, err error) {
	commenter, ok := cache.GetUser(userId)
	if !ok {
		commenter, err = db.FindUserInfo(userId)
		if err != nil {
			log.Panicln("service-GetCommenter: 获取评论者信息失败，", err)
			return model.User{}, err
		}
	}
	cache.SetUser(&commenter)
	return commenter, nil
}

func GetAuthor(videoID uint) (authorID uint, err error) {

	video, ok := cache.GetVideo(videoID)
	if !ok {
		video, err = db.GetVideo(videoID)
		if err != nil {
			log.Panicln("没有在数据库中查到该视频的相关信息", err)
			return
		}
	}
	cache.SetVideo(&video)
	authorID = video.AuthorID
	return authorID, nil
}

func CheckFollowing(hostID, guestID uint) bool {
	following, ok := cache.CheckFollowing(hostID, guestID)
	if !ok {
		following = db.CheckFollowing(hostID, guestID)
	}
	cache.SetFollowing(hostID, guestID, following)
	return following
}

// DeleteComment 删除评论
func DeleteComment(commentID uint, videoID uint) (err error) {
	//先查是否存在
	comment, err := db.GetComment(commentID)
	if err != nil {
		log.Panicf("数据库中没有id为%v的评论\n", commentID)
		//不存在, 就当作直接删除成功
		return nil
	}
	err = db.DelCommentTx(commentID, videoID)
	cache.DelComment(comment)
	if err != nil {
		return err
	}
	return nil
}

func CommentList(userId uint, videoID uint) (commentResponseList []common.CommentResp, err error) {
	var commentList []model.Comment
	commentList = cache.GetCommentList(videoID)
	if len(commentList) <= 0 {
		commentList, err = db.GetCommentList(videoID)
	}
	cache.SetMultiComment(&commentList)
	// log.Println("commentList: ", commentList)
	if err != nil {
		log.Println("service-CommentList: 查表获取评论列表时失败")
		return []common.CommentResp{{}}, nil
	}
	for i := 0; i < len(commentList); i++ {
		getUser, ok := cache.GetUser(commentList[i].UserID)
		if !ok {
			getUser, err = db.FindUserInfo(commentList[i].UserID)
			if err != nil {
				log.Println("无法找到评论者 ", getUser.ID, "，已略过此条评论 ", commentList[i].ID)
				continue
			}
		}
		cache.SetUser(&getUser)
		responseComment := common.CommentResp{
			ID:         commentList[i].ID,
			Content:    commentList[i].CommentText,
			CreateDate: commentList[i].CreatedAt.Format("01-02"), // mm-dd
			//这里应该是失误, 已更正位查询和userID的关注关系
			User: util.PackUserInfo(getUser, CheckFollowing(userId, commentList[i].UserID)),
		}
		commentResponseList = append(commentResponseList, responseComment)
	}
	return commentResponseList, nil
}
