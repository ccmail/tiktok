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
	"time"
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
	saveName := fmt.Sprint(parseToken.UserId, time.Now().UnixNano(), fileName)
	//去除掉一些会引起错误的非法字符z
	saveName = util.RemoveIllegalChar(saveName)
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
	playerUrl, err := middleware.OssUploadFromPath(saveName, savePath)
	if err != nil {
		log.Panicln("文件上传云服务器失败", err)
		return "", err
	}
	//更改文件后缀, 作为封面
	coverName := fmt.Sprint(saveName[:len(saveName)-4], ".jpeg")
	//获取第一帧作为封面
	cover, err := middleware.ExampleReadFrameAsJpeg(savePath, 1)
	if err != nil {
		log.Panicln("ffmpeg生成封面失败", err)
		return "", err
	}
	coverURL, err := middleware.OssUploadFromReader(coverName, cover)
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

// PublishList
// 这里逻辑存在一些问题, 是只有用户点开自己的发布列表使用这个方法, 还是点开其他用户的发布列表也使用这个方法?
// 目前按照点开其他用户的发布列表也是使用这个方法
// 等Feed完成之后需要测试
func (t *VideoService) PublishList(guestID uint, hostToken string) (resultList []common.VideoResp, err error) {
	guestInfo, err := mapper.FindUserInfo(guestID)
	if err != nil {
		log.Panicln("查找用户信息失败, 没有找到该用户的相关信息", err)
		return resultList, err
	}

	hostID := getHostIDFromToken(hostToken)

	author := packUserInfo(guestInfo, hostID)

	videoList, err := mapper.FindVideosByUserID(guestInfo.ID)
	if err != nil {
		log.Panicln("获取视频列表失败", err)
		return resultList, err
	}

	//需要展示的列表信息
	for i := 0; i < len(videoList); i++ {
		resultList = append(resultList, packVideoInfo(videoList[i], author, hostID))
	}
	return resultList, nil
}

func (t *VideoService) Feed(token string, strLastTime string) ([]common.VideoResp, int64, error) {

	hostID := getHostIDFromToken(token)
	var nextTime = time.Now()
	lastTime := time.Now()
	if strLastTime != "" {
		//没有传入last_time字段
		parse, err := time.Parse(strLastTime, "2006-01-02 15:03:04")
		if err != nil {
			log.Println("传入的latest_time格式有误", err)
		} else {
			lastTime = parse
		}
	}

	feedVideos := make([]common.VideoResp, 0, 30)

	videos, err := mapper.FindVideosByLastTime(lastTime)
	if err != nil {
		log.Panicln("根据时间请求视频时失败", err)
		return feedVideos, 0, err
	}
	usersID := make([]uint, len(videos))
	for i := range usersID {
		usersID[i] = videos[i].AuthorID
	}
	multiUser, err := mapper.FindMultiUserInfo(usersID)
	if err != nil {
		log.Panicln("获取作者信息失败", err)
		return feedVideos, 0, err
	}
	for _, video := range videos {
		if nextTime.After(video.CreatedAt) {
			nextTime = video.CreatedAt
		}
		if _, ok := multiUser[video.AuthorID]; !ok {
			log.Printf("没有查找到视频id为%v的作者信息, 作者信息应当为%v\n", video.ID, video.AuthorID)
			continue
		}
		userResp := packUserInfo(multiUser[video.AuthorID], hostID)
		feedVideos = append(feedVideos, packVideoInfo(video, userResp, hostID))
	}
	if len(feedVideos) == 0 {
		return feedVideos, 0, errors.New("没有获取到Feed信息")
	}
	return feedVideos, nextTime.Unix(), nil
}

func packVideoInfo(videos model.Video, author common.AuthorResp, hostID uint) common.VideoResp {
	returnVideo := common.VideoResp{
		Id:            videos.ID,
		Author:        author,
		PlayUrl:       videos.PlayUrl,
		CoverUrl:      videos.CoverUrl,
		FavoriteCount: videos.FavoriteCount,
		CommentCount:  videos.CommentCount,
		IsFavorite:    IsFavorite(hostID, videos.ID),
		Title:         videos.Title,
	}
	return returnVideo
}

// 包装user信息, 将其包装为response格式
func packUserInfo(guestInfo model.User, hostID uint) common.AuthorResp {
	author := common.AuthorResp{
		AuthorId:      guestInfo.ID,
		Name:          guestInfo.Name,
		FollowCount:   guestInfo.FollowCount,
		FollowerCount: guestInfo.FollowerCount,
		IsFollow:      IsFollowing(hostID, guestInfo.ID),
	}
	return author
}

// IsFavorite 查询某用户是否点赞某视频
func IsFavorite(uid uint, vid uint) bool {
	if uid == 0 {
		//uid为0代表用户未登录, 默认返回false, 代表未关注
		return false
	}
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
func getHostIDFromToken(hostToken string) uint {
	var hostID uint
	if hostToken != "" {
		hostInfo, err := middleware.ParseTokenCJS(hostToken)
		if err != nil {
			//就算token解析失败也不应该拒绝访问publishList
			log.Println("请求作品列表的用户的token解析失败", err)
		} else {
			hostID = hostInfo.UserId
		}
	}
	//log.Println(hostID)
	return hostID
}
