package controller

//
//import (
//	"fmt"
//	"github.com/gin-gonic/gin"
//	logging "github.com/sirupsen/logrus"
//	"gorm.io/gorm"
//	"net/http"
//	"os"
//	"path/filepath"
//	"strconv"
//	"strings"
//	"tiktok/model"
//	"tiktok/pkg/common"
//	middleware "tiktok/pkg/middleware"
//	"tiktok/service"
//	"time"
//)
//
//// Publish 投稿控制层
//func Publish(c *gin.Context) {
//	// 中间件验证token后，获取userId
//	rawUserId, _ := c.Get("user_id")
//	var userId uint
//	if v, ok := rawUserId.(uint); ok {
//		userId = v
//	}
//	// 接收请求参数信息
//	title := c.PostForm("title")
//	data, err := c.FormFile("data")
//	if err != nil {
//		c.JSON(http.StatusOK, common.BaseResponse{
//			StatusCode: 1,
//			StatusMsg:  err.Error(),
//		})
//		return
//	}
//
//	// 返回至前端页面的展示信息
//	filename := filepath.Base(data.Filename)
//	finalName := fmt.Sprintf("%d_%s", userId, filename)
//
//	//从这里开始重写
//	// 先存储到本地文件夹，再保存到云端，获取封面后删除本地文件
//	savePath := filepath.Join("./videos/", finalName)
//	err = c.SaveUploadedFile(data, savePath)
//	if err != nil {
//		c.JSON(http.StatusOK, common.BaseResponse{
//			StatusCode: 1,
//			StatusMsg:  err.Error(),
//		})
//		return
//	}
//
//	// 从本地上传到云端，并获取云端地址
//	playUrl, err := service.OssUploadFromPath(finalName, savePath)
//	if err != nil {
//		c.JSON(http.StatusOK, common.BaseResponse{
//			StatusCode: 1,
//			StatusMsg:  err.Error(),
//		})
//		return
//	}
//
//	// 直接传至云端，不用存储到本地
//	coverName := strings.Replace(finalName, ".mp4", ".jpeg", 1)
//	img, _ := service.ExampleReadFrameAsJpeg(savePath, 1) // 获取第1帧作为封面
//
//	coverUrl, err := service.OssUploadFromReader(coverName, img)
//	if err != nil {
//		c.JSON(http.StatusOK, common.BaseResponse{
//			StatusCode: 1,
//			StatusMsg:  err.Error(),
//		})
//		return
//	}
//
//	// 删除保存在本地中的视频
//	err = os.Remove(savePath) // ignore_security_alert
//	if err != nil {
//		logging.Info(err)
//	}
//
//	// 保存发布信息至数据库,刚开始发布，喜爱和评论默认为0
//	video := model.Video{
//		Model:         gorm.Model{},
//		AuthorID:      userId,
//		PlayUrl:       playUrl,
//		CoverUrl:      coverUrl,
//		FavoriteCount: 0,
//		CommentCount:  0,
//		Title:         title,
//	}
//	service.CreateVideo(&video)
//
//	c.JSON(http.StatusOK, common.BaseResponse{
//		StatusCode: 0,
//		StatusMsg:  finalName + "--uploaded successfully",
//	})
//}
//
//func PublishList(c *gin.Context) {
//	// 中间件鉴权 token 后，获取 hostID （即发起请求的用户 id ）
//	getHostId, _ := c.Get("user_id")
//	var HostId uint
//	if v, ok := getHostId.(uint); ok {
//		HostId = v
//	}
//	// fmt.Println("hostID: ", HostId)
//
//	// 查询 guestID （即被查询发布列表的用户 id ）用户的所有视频，返回页面
//	getGuestId := c.Query("user_id")
//	id, _ := strconv.Atoi(getGuestId)
//	GuestId := uint(id)
//	// fmt.Println("guestID: ", GuestId)
//
//	// 根据传入 id 查找用户
//	getUser, err := service.GetUser(GuestId)
//	if err != nil {
//		c.JSON(http.StatusOK, common.BaseResponse{
//			StatusCode: 1,
//			StatusMsg:  "Target user not found.",
//		})
//		c.Abort()
//		return
//	}
//
//	returnAuthor := common.ReturnAuthor{
//		AuthorId:      getUser.ID,
//		Name:          getUser.Name,
//		FollowCount:   getUser.FollowCount,
//		FollowerCount: getUser.FollowerCount,
//		IsFollow:      service.IsFollowing(HostId, GuestId),
//	}
//
//	// 根据用户id查找 所有相关视频信息
//	videoList := service.GetVideoList(GuestId)
//	if len(videoList) == 0 {
//		c.JSON(http.StatusOK, common.VideoListResponse{
//			BaseResponse: common.BaseResponse{
//				StatusCode: 1,
//				StatusMsg:  "Empty video list",
//			},
//			VideoList: nil,
//		})
//		return
//	}
//
//	//需要展示的列表信息
//	var returnVideoList []common.ReturnVideo
//	for i := 0; i < len(videoList); i++ {
//		returnVideo := common.ReturnVideo{
//			VideoId:       videoList[i].ID,
//			Author:        returnAuthor,
//			PlayUrl:       videoList[i].PlayUrl,
//			CoverUrl:      videoList[i].CoverUrl,
//			FavoriteCount: videoList[i].FavoriteCount,
//			CommentCount:  videoList[i].CommentCount,
//			IsFavorite:    service.IsFavorite(HostId, videoList[i].ID),
//			Title:         videoList[i].Title,
//		}
//		returnVideoList = append(returnVideoList, returnVideo)
//	}
//	c.JSON(http.StatusOK, common.VideoListResponse{
//		BaseResponse: common.BaseResponse{
//			StatusCode: 0,
//			StatusMsg:  "Get publish list: success",
//		},
//		VideoList: returnVideoList,
//	})
//
//}
//
//func Feed(c *gin.Context) {
//	strToken := c.Query("token")
//	var haveToken bool
//	if strToken == "" {
//		haveToken = false
//	} else {
//		haveToken = true
//	}
//	strLastTime := c.Query("latest_time")
//
//	lastTime, err := strconv.ParseInt(strLastTime, 10, 32)
//	if err != nil { // 没传 latest_time 参数的情况
//		lastTime = time.Now().Unix()
//	}
//
//	feedVideoList := make([]common.FeedVideo, 0)
//	videoList, _ := service.FeedSerivce(lastTime)
//	// fmt.Println("VideoList length: ", len(videoList))
//	var newTime int64 // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
//
//	// 逐视频赋予参数
//	for _, v := range videoList {
//		var tmp common.FeedVideo
//		tmp.Id = v.ID
//		tmp.PlayUrl = v.PlayUrl
//		//tmp.Author = //依靠用户信息接口查询
//		user, err := service.GetUser(v.AuthorID)
//		var feedUser common.FeedUser
//		if err == nil { // 用户存在
//			feedUser.Id = user.ID
//			feedUser.FollowerCount = user.FollowerCount
//			feedUser.FollowCount = user.FollowCount
//			feedUser.Name = user.Name
//			//add
//			// feedUser.TotalFavorited = user.TotalFavorited
//			// feedUser.FavoriteCount = user.FavoriteCount
//			feedUser.IsFollow = false
//			if haveToken {
//				// 查询是否关注
//				tokenStruct, ok := middleware.ParseToken(strToken)
//
//				// token 超时
//				if time.Now().Unix() > tokenStruct.ExpiresAt {
//					c.JSON(http.StatusOK, common.BaseResponse{
//						StatusCode: 402,
//						StatusMsg:  "Expired token",
//					})
//					return
//				}
//				// token 合法
//				if ok {
//					uid1 := tokenStruct.UserId // 用户id
//					uid2 := v.AuthorID         // 视频发布者id
//					if service.IsFollowing(uid1, uid2) {
//						feedUser.IsFollow = true
//					}
//
//					vid := v.ID                        // 视频id
//					if service.IsFavorite(uid1, vid) { //有点赞记录
//						tmp.IsFavorite = true
//					}
//				}
//			}
//		}
//		tmp.Author = feedUser
//		tmp.CommentCount = v.CommentCount
//		tmp.CoverUrl = v.CoverUrl
//		tmp.FavoriteCount = v.FavoriteCount
//		tmp.IsFavorite = false
//		tmp.Title = v.Title
//		feedVideoList = append(feedVideoList, tmp)
//		newTime = v.CreatedAt.Unix()
//	}
//	// fmt.Println("feedVideoList length: ", len(feedVideoList))
//	if len(feedVideoList) > 0 {
//		c.JSON(http.StatusOK, common.FeedResponse{
//			BaseResponse: common.BaseResponse{
//				StatusCode: 0,
//				StatusMsg:  "Get feed video list: success"},
//			VideoList: feedVideoList,
//			NextTime:  uint(newTime),
//		})
//	} else {
//		c.JSON(http.StatusOK, common.FeedNoVideoResponse{
//			BaseResponse: common.BaseResponse{StatusCode: 0},
//			NextTime:     0, //重新循环
//		})
//	}
//
//}
