package service

import (
	"log"
	"tiktok/mapper"
	"tiktok/model"
	"tiktok/pkg/common"
)

// PostCommentService 发布评论
func PostCommentService(userId uint, text string, videoId uint) (model.Comment, error) {
	newComment := model.Comment{
		VideoID:     videoId,
		UserID:      userId,
		CommentText: text,
		Valid:       true,
	}

	err := mapper.NewCommentTx(newComment)
	if err != nil {
		return model.Comment{}, err
	}
	return newComment, nil
}

func GetCommenter(userId uint) (model.User, error) {
	commenter, err := mapper.FindUserInfo(userId)
	if err != nil {
		log.Panicln("service-GetCommenter: 获取评论者信息失败，", err)
		return model.User{}, err
	}
	return commenter, nil
}

func GetAuthor(videoID uint) (uint, error) {
	authorID, err := mapper.GetVideoAuthor(videoID)
	if err != nil {
		log.Println("service-GetAuthor: 获取作者失败，", err)
		return 0, err
	}
	return authorID, nil
}

// DeleteCommentService 删除评论
func DeleteCommentService(commentID uint, videoID uint) error {
	err := mapper.DelCommentTx(commentID, videoID)
	if err != nil {
		return err
	}
	return nil
}

func CommentListService(userId uint, videoID uint) (commentRespionseList []common.CommentResponse, err error) {
	commentList, err := mapper.GetCommentList(videoID)
	// log.Println("commentList: ", commentList)
	if err != nil {
		log.Println("service-CommentListService: 查表获取评论列表时失败")
		return []common.CommentResponse{{}}, nil
	}
	for i := 0; i < len(commentList); i++ {
		getUser, err := mapper.FindUserInfo(commentList[i].UserID)

		if err != nil {
			log.Println("无法找到评论者 ", getUser.ID, "，已略过此条评论 ", commentList[i].ID)
			continue
		}
		responseComment := common.CommentResponse{
			ID:         commentList[i].ID,
			Content:    commentList[i].CommentText,
			CreateDate: commentList[i].CreatedAt.Format("01-02"), // mm-dd
			User: common.CommenterInfo{
				UserID:        getUser.ID,
				Username:      getUser.Name,
				FollowCount:   getUser.FollowCount,
				FollowerCount: getUser.FollowerCount,
				IsFollow:      mapper.CheckFollowing(userId, commentList[i].ID),
			},
		}
		commentRespionseList = append(commentRespionseList, responseComment)
	}
	return commentRespionseList, nil
}
