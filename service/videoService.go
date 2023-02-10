package service

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"tiktok/mapper"
	"tiktok/model"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	"gorm.io/gorm"
)

// func PublishService(userId uint, data *multipart.FileHeader) error {

// }

// OssUpload 上传至云端Oss，返回url
func OssUploadFromPath(filename string, filepath string) (url string, err error) {
	err = mapper.Bucket.PutObjectFromFile("short_video/"+filename, filepath)
	if err != nil {
		return "", err
	}
	url = "https://" + mapper.BucketName + "." + mapper.EndPoint + "/short_video/" + filename
	return url, nil
}

func OssUploadFromReader(filename string, data io.Reader) (url string, err error) {
	err = mapper.Bucket.PutObject("short_video/"+filename, data)
	if err != nil {
		return "", err
	}
	url = "https://" + mapper.BucketName + "." + mapper.EndPoint + "/short_video/" + filename
	return url, nil
}

// CreateVideo 添加一条视频信息
func CreateVideo(video *model.Video) {
	mapper.DBConn.Table("videos").Create(&video)
}

// ExampleReadFrameAsJpeg 获取封面
func ExampleReadFrameAsJpeg(inFileName string, frameNum int) io.Reader {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		panic(err)
	}
	return buf
}

// IsFavorite 查询某用户是否点赞某视频
func IsFavorite(uid uint, vid uint) bool {
	var total int64
	err := mapper.DBConn.Table("likes").Where("user_id = ? AND video_id = ? AND is_like = ?", uid, vid, true).Count(&total).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	if total == 0 {
		return false
	}
	return true
}

const maxVideoNum = 30 // feed每次返回的最大视频数量

// FeedSerivce 获得视频列表
func FeedSerivce(lastTime int64) ([]model.Video, error) {
	//t := time.Now()
	//fmt.Println(t)
	strTime := time.Unix(lastTime, 0).Format("2006-01-02 15:03:04")

	// fmt.Println("查询的时间", strTime)
	VideoList := make([]model.Video, 0)
	err := mapper.DBConn.Table("videos").Where("created_at < ?", strTime).Order("created_at desc").Limit(maxVideoNum).Find(&VideoList).Error
	return VideoList, err
}

// AddCommentCount add comment_count
func AddCommentCount(videoId uint) error {

	if err := mapper.DBConn.Table("videos").Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count + 1")).Error; err != nil {
		return err
	}
	return nil
}

// ReduceCommentCount reduce comment_count
func ReduceCommentCount(videoId uint) error {

	if err := mapper.DBConn.Table("videos").Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count - 1")).Error; err != nil {
		return err
	}
	return nil
}

// GetVideoAuthor get video author
func GetVideoAuthor(videoId uint) (uint, error) {
	var video model.Video
	if err := mapper.DBConn.Table("videos").Where("id = ?", videoId).Find(&video).Error; err != nil {
		return video.ID, err
	}
	return video.AuthorID, nil
}

// GetVideoList 根据用户id查找 所有与该用户相关视频信息
func GetVideoList(userId uint) []model.Video {
	var videoList []model.Video
	mapper.DBConn.Table("videos").Where("author_id=?", userId).Find(&videoList)
	return videoList
}
