package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"tiktok/config"
	"tiktok/mapper"
	"tiktok/model"
	"tiktok/pkg/common"
	"tiktok/pkg/middleware"
	util2 "tiktok/pkg/util"
	"time"
)

// Publish
//上传投稿视频, 先对token鉴权, 后从token中获取用户id, 再将视频重命名, 上传到文件服务器, 同时使用ffmpeg生成封面, 一并上传到文件服务器
// 如果上述流程中出现了一场, 返回值为("",err), err为对应的错误信息
// 流程正常, 返回上传到文件服务器上的名称和nil, 即(saveName,nil)
func (t *VideoService) Publish(file *multipart.FileHeader, title, token string) (string, error) {
	parseToken, err := middleware.ParseTokenCJS(token)
	if err != nil {
		return "", err
	}
	fileName := filepath.Base(file.Filename)
	//存储到文件服务器中的文件名字, 为了避免重复, 这里使用用户id+时间戳+文件名命名
	saveName := fmt.Sprint(parseToken.UserId, time.Now().UnixNano(), fileName)
	//去除掉一些会引起错误的非法字符z
	saveName = util2.RemoveIllegalChar(saveName)
	savePath := filepath.Join("./videos/", saveName)
	//保存视频到本地
	err = util2.SaveFileLocal(file, savePath)
	if err != nil {
		log.Panicln("文件保存失败", err)
		return "", err
	}
	//最后删除本地视频
	defer util2.RemoveFileLocal(savePath)

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

	hostID := util2.GetHostIDFromToken(hostToken)

	author := util2.PackUserInfo(guestInfo, mapper.CheckFollowing(hostID, guestInfo.ID))

	videoList, err := mapper.FindVideosByUserID(guestInfo.ID)
	if err != nil {
		log.Panicln("获取视频列表失败", err)
		return resultList, err
	}

	//需要展示的列表信息
	for i := 0; i < len(videoList); i++ {
		resultList = append(resultList, util2.PackVideoInfo(videoList[i], author, mapper.IsFavorite(hostID, videoList[i].ID)))
	}
	return resultList, nil
}

func (t *VideoService) Feed(token string, strLastTime string) (vResp []common.VideoResp, nTime int64, err error) {

	hostID := util2.GetHostIDFromToken(token)
	var nextTime = time.Now()
	lastTime := time.Now()
	if strLastTime != "" && strLastTime != "0" && len(strLastTime) == 10 {
		i, err := strconv.ParseInt(strLastTime, 10, 64)
		if err != nil {
			log.Println("传入字符串格式有误")
		} else {
			lastTime = time.Unix(i, 0)
		}
	}
	//查出来之后需要查询用户的点赞信息以及用户和人家的关注信息
	videoInfos := mapper.GetFeedCache(lastTime)

	//缓存中的feed数量不够, 查mysql
	if len(videoInfos) < config.MaxFeedVideoCount {
		videoInfos, err = mapper.FindVideosByLastTime(lastTime)
		if err != nil {
			log.Panicln("根据时间向Mysql请求视频时失败", err)
			return vResp, nTime, err
		}
		//	这里还需要将videos写入cache, 直接全部写入, 以达到将存在值的cache存活时间更新
		mapper.SetMultiFeedCache(&videoInfos)
	}
	if len(videoInfos) == 0 {
		return vResp, lastTime.Unix(), errors.New("没有更多视频, 等会试试吧")
	}

	//查缓存, 点赞信息的缓存
	isFavorite, likeNoCache := mapper.CheckMultiFavoriteCache(hostID, &videoInfos)
	if len(likeNoCache) > 0 {
		//	有几个没查到缓存, 需要查数据库
		err := mapper.CheckLikesNoHit(hostID, &isFavorite, &likeNoCache)
		if err != nil {
			log.Println("数据库中也不存在对应的点赞关系, 不用管了, 直接当作不存在点赞信息")
		}
		//每次查都需要更新信息在cache中的存活时间, 查到之后要写到cache中
		mapper.SetMultiFavoriteCache(hostID, &videoInfos, &isFavorite)
	}

	//设置up的id数组
	guestIDs := make([]uint, len(videoInfos))
	for i := range guestIDs {
		guestIDs[i] = videoInfos[i].AuthorID
	}

	//关注信息查缓存
	isFollow, followNoCache := mapper.CheckMultiFollowingCache(hostID, &guestIDs)
	if len(followNoCache) > 0 {
		err := mapper.CheckMultiFollowNoHit(hostID, &isFollow, &followNoCache)
		if err != nil {
			log.Println("数据库中也不存在对应的关注关系, 不用管了, 直接当作不存在关注信息")
		}
		//	之后将查到的关注信息写入到redis中
		mapper.SetMultiFollowingCache(hostID, &guestIDs, &isFollow)
	}

	//作者信息查缓存
	userInfo, userNoCache := mapper.GetMultiUserCache(guestIDs)
	if len(userNoCache) > 0 {
		err := mapper.GetMultiUserInfoNoHit(&userInfo, &userNoCache)
		if err != nil {
			log.Println("数据库中也不存在这个user 可能出错了")
		}
		//	之后将查到的信息写入到redis中
		mapper.SetMultiUserCache(&userInfo)
	}

	//按倒叙传入的视频, 最后的时间最前
	nextTime = videoInfos[len(videoInfos)-1].CreatedAt
	vResp = util2.PackVideoListInfo(videoInfos, userInfo, isFollow, isFavorite)

	if len(vResp) == 0 {
		return vResp, 0, errors.New("没有获取到Feed信息")
	}
	return vResp, nextTime.Unix(), nil
}
