package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"mime/multipart"
	"path/filepath"
	"tiktok/mapper"
	"tiktok/model"
	"tiktok/pkg/common"
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
		log.Panicln("文件保存失败", err)
		return "", err
	}
	//最后删除本地视频
	defer util.RemoveFileLocal(savePath)

	//从本地上传到云服务器
	playerUrl, err := OssUploadFromPath(saveName, savePath)
	if err != nil {
		log.Panicln("文件上传云服务器失败", err)
		return "", err
	}
	//更改文件后缀, 作为封面
	coverName := fmt.Sprint(saveName[:len(saveName)-4], ".jpeg")
	//获取第一帧作为封面
	cover, err := ExampleReadFrameAsJpeg(savePath, 1)
	if err != nil {
		log.Panicln("ffmpeg生成封面失败", err)
		return "", err
	}
	coverURL, err := OssUploadFromReader(coverName, cover)
	if err != nil {
		log.Panicln("视频封面上传失败", err)
		return "", err
	}
	video := model.Video{
		Model:    gorm.Model{},
		AuthorID: parseToken.UserId,
		PlayUrl:  playerUrl,
		CoverUrl: coverURL,
		Title:    title,
	}
	//service层调用dao(mapper)层
	err = mapper.CreateVideo(&video)
	if err != nil {
		log.Panicln("视频信息插入数据库失败", err)
		return "", err
	}
	return saveName, nil
}

func (t *VideoService) PublishList(guestID uint, hostToken string) (resultList []common.ReturnVideo, err error) {
	guestInfo, err := mapper.FindUserInfo(guestID)
	if err != nil {
		log.Panicln("查找用户信息失败, 没有找到该用户的相关信息", err)
		return resultList, err
	}
	hostInfo, err := middleware.ParseTokenCJS(hostToken)
	if err != nil {
		log.Panicln("请求作品列表的用户的token解析失败", err)
		return resultList, err
	}

	author := common.ReturnAuthor{
		AuthorId:      guestInfo.ID,
		Name:          guestInfo.Name,
		FollowCount:   guestInfo.FollowCount,
		FollowerCount: guestInfo.FollowerCount,
		IsFollow:      IsFollowing(hostInfo.UserId, guestInfo.ID),
	}

	videoList, err := mapper.FindVideoList(guestInfo.ID)
	if err != nil {
		log.Panicln("获取视频列表失败", err)
		return resultList, err
	}

	//需要展示的列表信息
	for i := 0; i < len(videoList); i++ {
		returnVideo := common.ReturnVideo{
			VideoId:       videoList[i].ID,
			Author:        author,
			PlayUrl:       videoList[i].PlayUrl,
			CoverUrl:      videoList[i].CoverUrl,
			FavoriteCount: videoList[i].FavoriteCount,
			CommentCount:  videoList[i].CommentCount,
			IsFavorite:    IsFavorite(hostInfo.UserId, videoList[i].ID),
			Title:         videoList[i].Title,
		}
		resultList = append(resultList, returnVideo)
	}
	return resultList, nil
}
