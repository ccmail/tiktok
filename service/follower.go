package service

import (
	"errors"
	"log"
	"tiktok/mapper/cache"
	"tiktok/mapper/db"
	"tiktok/pkg/common"
	"tiktok/pkg/middleware"
	"tiktok/pkg/util"
)

// Follow
// 1.向follow表中添加host和guest的关注关系
// 2.更改User表中对应的粉丝/关注up数量
func Follow(token string, guestID uint, isConcern bool) error {
	parseToken, err := middleware.ParseTokenCJS(token)
	if err != nil {
		return err
	}
	record, exist := db.ExistFollowRecord(parseToken.UserId, guestID)
	//如果不存在的话, 向follow插入信息, 并更新user表的字段
	if !exist {
		err := db.TranCreateFollow(parseToken.UserId, guestID, isConcern)
		if err != nil {
			log.Panicln("插入关注关系到数据库时发生错误")
			return err
		}
		cache.SetFollowing(parseToken.UserId, guestID, isConcern)
		cache.DelUser(guestID, parseToken.UserId)
		return nil
	}

	if record.IsFollow == isConcern || (!record.IsFollow && !isConcern) {
		//关注关系没有变化, 或者已经处于未关注状态
		return nil
	}
	//更新host的关注数, 并增加host关注的up数量, 增加up的粉丝数量

	err = db.TranUpdateFollow(parseToken.UserId, guestID, isConcern)
	if err != nil {
		log.Panicln("更新关注关系到数据库时发生错误")
		return err
	}
	cache.SetFollowing(parseToken.UserId, guestID, isConcern)
	cache.DelUser(guestID, parseToken.UserId)
	return nil
}

// FollowList 用户token去请求guestID对应的关注列表
func FollowList(token string, guestID uint) (resList []common.UserInfoResp, err error) {
	hostID := util.GetHostIDFromToken(token)

	//这里请求出来有两步, 如果tokenID==guestID, 直接返回, 默认全部关注
	//如果tokenID!=guestID, 还需要判断请求出来的user和token的关注关系
	guestIDList := cache.GetFollowIDList(guestID)
	if len(guestIDList) == 0 {
		guestIDList, err = db.FindMultiConcern(guestID)
		if err != nil {
			log.Panicln("从数据库中查找关注列表信息时失败")
			return resList, err
		}
	}
	cache.SetFollowIDList(guestID, &guestIDList)

	userList, uid := cache.GetMultiUser(&guestIDList)
	if len(uid) > 0 {
		err := db.GetMultiUserInfoNoHit(&userList, &uid)
		if err != nil {
			log.Panicln("有部分用户信息查找失败了")
			return resList, err
		}
	}
	cache.SetMultiUser(&userList)

	followList, bad := cache.CheckMultiFollowing(hostID, &guestIDList)
	if len(bad) > 0 {
		err := db.CheckMultiFollowNoHit(hostID, &followList, &bad)
		if err != nil {
			log.Println("有部分关注信息查找不到")
			return resList, err
		}
	}
	cache.SetMultiFollowing(hostID, &guestIDList, &followList)

	resList = make([]common.UserInfoResp, 0, len(userList))
	for i, upInfo := range userList {
		isFollowing := true
		if hostID != guestID {
			isFollowing = followList[i]
		}
		resList = append(resList, util.PackUserInfo(upInfo, isFollowing))
	}
	return resList, nil
}

// FollowerList 请求粉丝列表, 和请求关注列表逻辑一致, 不同的是需要更改一下查询数据库的信息
func FollowerList(token string, guestID uint) (resList []common.UserInfoResp, err error) {
	hostID := util.GetHostIDFromToken(token)

	userIDList := cache.GetFansIDList(guestID)
	if len(userIDList) == 0 {
		userIDList, err = db.FindMultiFollower(guestID)
		if err != nil {
			log.Panicln("查找粉丝列表信息时失败")
			return resList, err
		}
	}
	cache.SetFansIDList(guestID, &userIDList)

	//获取到guestID关注的up的用户信息
	userInfoList, uid := cache.GetMultiUser(&userIDList)
	if len(uid) > 0 {
		err := db.GetMultiUserInfoNoHit(&userInfoList, &uid)
		if err != nil {
			log.Panicln("查找粉丝个人信息时失败")
			return resList, err
		}
	}

	followList, bad := cache.CheckMultiFollowing(hostID, &userIDList)
	if len(bad) > 0 {
		err := db.CheckMultiFollowNoHit(hostID, &followList, &bad)
		if err != nil {
			log.Println("有部分关注信息查找不到")
			return resList, err
		}
	}
	cache.SetMultiFollowing(hostID, &userIDList, &followList)

	resList = make([]common.UserInfoResp, 0, len(userInfoList))
	for i, upInfo := range userInfoList {
		isFollowing := true
		if hostID != guestID {
			isFollowing = followList[i]
		}
		resList = append(resList, util.PackUserInfo(upInfo, isFollowing))
	}
	return resList, nil
}

// FriendList	这里要做一下过滤, 只有自己能看自己的好友列表, 当guestID与token不符时应当直接返回
func FriendList(token string, guestID uint) (resList []common.FriendInfo, err error) {
	hostID := util.GetHostIDFromToken(token)

	if hostID != guestID {
		return resList, errors.New("请求好友列表的id和token不一致, 不允许偷看别人的好友列表")
	}

	ups := cache.GetFollowIDList(guestID)
	if len(ups) == 0 {
		ups, err = db.FindMultiConcern(guestID)
		if err != nil {
			log.Panicln("从数据库中查找关注列表信息时失败")
			return resList, err
		}
	}
	cache.SetFollowIDList(guestID, &ups)

	fans := cache.GetFansIDList(guestID)
	if len(fans) == 0 {
		fans, err = db.FindMultiFollower(guestID)
		if err != nil {
			log.Panicln("查找粉丝列表信息时失败")
			return resList, err
		}
	}
	cache.SetFansIDList(guestID, &fans)

	mp := make(map[uint]interface{})
	for i := 0; i < len(fans); i++ {
		mp[fans[i]] = struct{}{}
	}
	friends := make([]uint, 0, len(fans)>>1)
	for i := 0; i < len(ups); i++ {
		if _, ok := mp[ups[i]]; ok {
			friends = append(friends, ups[i])
		}
	}

	//获取到guestID关注的up的用户信息
	userInfoList, uid := cache.GetMultiUser(&friends)
	if len(uid) > 0 {
		err := db.GetMultiUserInfoNoHit(&userInfoList, &uid)
		if err != nil {
			log.Panicln("获取好友信息时失败了")
			return resList, err
		}
	}
	resList = make([]common.FriendInfo, 0, len(userInfoList))

	//userResp := make([]common.UserInfoResp, 0, len(userInfoList))
	for _, userInfo := range userInfoList {
		//查找聊天记录, 返回时间最后的
		send := db.GetSendMessage(hostID, userInfo.ID)
		receive := db.GetReceiveMessage(hostID, userInfo.ID)
		var msgType = int64(0)
		mes := receive.MessageText
		if send.CreatedAt.After(receive.CreatedAt) {
			msgType = 1
			mes = send.MessageText
		}

		resList = append(resList, common.FriendInfo{
			Message:      mes,
			MsgType:      msgType,
			UserInfoResp: util.PackUserInfo(userInfo, true),
		})
	}
	return resList, nil
}
