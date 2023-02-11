package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"mime/multipart"
	"path/filepath"
	"tiktok/model"
	"tiktok/pkg/middleware"
	"tiktok/util"
)

// Publish
//上传投稿视频, 先对token鉴权, 后从token中获取用户id, 再将视频重命名, 上传到文件服务器, 同时使用ffmpeg生成封面, 一并上传到文件服务器
// 如果上述流程中出现了一场, 返回值为("",err), err为对应的错误信息
// 流程正常, 返回上传到文件服务器上的名称和nil, 即(saveName,nil)
func (t *VideoService) Publish(file *multipart.FileHeader, title, token string) (string, error) {
	parseToken, ok := middleware.ParseToken(token)
	if !ok {
		return "", errors.New("token失效")
	}
	fileName := filepath.Base(file.Filename)
	//存储到文件服务器中的文件名字, 为了避免重复, 这里使用用户id+时间戳+文件名命名
	saveName := fmt.Sprint(parseToken.UserId, fileName)
	//saveName := fmt.Sprint(parseToken.UserId, time.Now().UnixNano(), fileName)
	savePath := filepath.Join("./videos/", saveName)
	//保存视频到本地
	err := util.SaveFileLocal(file, savePath)
	if err != nil {
		log.Println("文件保存失败", err)
		return "", errors.New("文件保存失败")
	}
	//最后删除本地视频
	defer util.RemoveFileLocal(savePath)

	//从本地上传到云服务器
	playerUrl, err := OssUploadFromPath(saveName, savePath)
	if err != nil {
		log.Println("文件上传云服务器失败", err)
		return "", errors.New("文件上传云服务器失败")
	}
	//更改文件后缀, 作为封面
	coverName := fmt.Sprint(saveName[:len(saveName)-4], ".jpeg")
	//获取第一帧作为封面
	cover, err := ExampleReadFrameAsJpeg(savePath, 1)
	if err != nil {
		log.Println("ffmpeg生成封面失败", err)
		return "", errors.New("ffmpeg生成封面失败")
	}
	coverURL, err := OssUploadFromReader(coverName, cover)
	if err != nil {
		log.Println("视频封面上传失败", err)
		return "", errors.New("视频封面上传失败")
	}
	video := model.Video{
		Model:    gorm.Model{},
		AuthorID: parseToken.UserId,
		PlayUrl:  playerUrl,
		CoverUrl: coverURL,
		Title:    title,
	}
	CreateVideo(&video)
	return saveName, nil
}
