package db

import (
	"gorm.io/gorm"
	"log"
	"tiktok/mapper"
	"tiktok/model"
)

func CreateComment(comment model.Comment) error {
	err := mapper.DBConn.Table("comments").Create(&comment).Error
	if err != nil {
		log.Println("mapper-CreateComment: 新建评论记录失败")
		return err
	}
	return nil
}

func NewCommentTx(newComment model.Comment) error {

	err1 := mapper.DBConn.Transaction(func(db *gorm.DB) error {
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

func DelComment(commentId uint) error {
	err := mapper.DBConn.Table("comments").Where("id = ?", commentId).Update("valid", false).Error
	if err != nil {
		log.Println("mapper-DelComment: 删除评论操作失败，", err)
		return err
	}
	return nil
}

func GetComment(cid uint) (ans model.Comment, err error) {
	err = mapper.DBConn.Model(&model.Comment{}).Where("id = ?", cid).Find(&ans).Error
	if err != nil {
		log.Panicf("没有查到id为%v的评论\n, 错误信息为:%v", cid, err)
		return
	}
	return ans, nil
}
func DelCommentTx(commentID uint, videoID uint) error {

	err1 := mapper.DBConn.Transaction(func(db *gorm.DB) error {
		if err := DelComment(commentID); err != nil {
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
	err = mapper.DBConn.Table("comments").Where("video_id=? AND valid=?", videoID, true).Order("created_at desc").Find(&commentList).Error
	if err != nil {
		log.Println("mapper-GetCommentList: 查表获取评论列表失败，", err)
		return []model.Comment{{}}, err
	}
	return commentList, nil
}
