package mapper

import (
	"log"
	"tiktok/model"

	"gorm.io/gorm"
)

func CreateComment(comment model.Comment) error {
	err := DBConn.Table("comments").Create(&comment).Error
	if err != nil {
		log.Println("mapper-CreateComment: 新建评论记录失败")
		return err
	}
	return nil
}

func NewCommentTx(newComment model.Comment) error {

	err1 := DBConn.Transaction(func(db *gorm.DB) error {
		if err := CreateComment(newComment); err != nil {
			return err
		}
		if err := AddCommentCount(newComment.VideoID); err != nil {
			return err
		}
		return nil
	})

	if err1 != nil {
		log.Println("mapper-NewCommentTx: 发表评论操作失败，", err1)
	}

	return nil
}

func DeleteComment(commentId uint) error {
	err := DBConn.Table("comments").Where("id = ?", commentId).Update("valid", false).Error
	if err != nil {
		log.Println("mapper-DeleteComment: 删除评论操作失败，", err)
		return err
	}
	return nil
}

func DelCommentTx(commentID uint, videoID uint) error {
	err1 := DBConn.Transaction(func(db *gorm.DB) error {
		if err := DeleteComment(commentID); err != nil {
			return err
		}
		if err := ReduceCommentCount(videoID); err != nil {
			return err
		}
		return nil
	})
	if err1 != nil {
		log.Println("mapper-DelCommentTx: 删除评论操作失败，", err1)
	}
	return nil
}

// GetCommentList 获取一个视频的评论列表
func GetCommentList(videoID uint) (commentList []model.Comment, err error) {
	err = DBConn.Table("comments").Where("video_id=? AND valid=?", videoID, true).Find(&commentList).Error
	if err != nil {
		log.Println("mapper-GetCommentList: 查表获取评论列表失败，", err)
		return []model.Comment{{}}, err
	}
	return commentList, nil
}
