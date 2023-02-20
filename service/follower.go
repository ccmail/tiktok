package service

import (
	"errors"
	"log"
	"tiktok/mapper/gorm"
	"tiktok/pkg/common"
	"tiktok/pkg/middleware"
	util2 "tiktok/pkg/util"
)

// Follow
// 1.向follow表中添加host和guest的关注关系
// 2.更改User表中对应的粉丝/关注up数量
func Follow(token string, guestID uint, isConcern bool) error {
	parseToken, err := middleware.ParseTokenCJS(token)
	if err != nil {
		return err
	}
	if parseToken.UserId == guestID {
		return errors.New("请不要自己关注自己")
	}
	record, exist := gorm.ExistFollowRecord(parseToken.UserId, guestID)
	//如果不存在的话, 向follow插入信息, 并更新user表的字段
	if !exist {
		err := gorm.CreatFollowRecord(parseToken.UserId, guestID, isConcern)
		if err != nil {
			log.Panicln("插入数据库时发生错误")
			return err
		}
		//从未关注过的话直接插, 更新完毕之后return掉
		err = gorm.UpdateUserFollowCount(parseToken.UserId, guestID, isConcern)
		if err != nil {
			log.Panicln("更新粉丝/关注数量时出错")
			return err
		}
		return nil
		//mapper.UpdateUserFollowCount(parseToken.UserId, guestID)
	}

	if record.IsFollow == isConcern || (!record.IsFollow && !isConcern) {
		//关注关系没有变化, 或者已经处于未关注状态
		return nil
	}
	//更新host的关注数, 并增加host关注的up数量, 增加up的粉丝数量
	err = gorm.UpdateUserFollowCount(parseToken.UserId, guestID, isConcern)
	if err != nil {
		log.Panicln("更新粉丝/关注数量时出错")
		return err
	}
	//向follow表中添加follow关系记录
	err = gorm.UpdateFollowRecord(parseToken.UserId, guestID, isConcern)
	if err != nil {
		log.Panicln("更新关注关系时失败")
		return err
	}
	return nil
}

// FollowList 用户token去请求guestID对应的关注列表
func FollowList(token string, guestID uint) (resList []common.UserInfoResp, err error) {
	hostID := util2.GetHostIDFromToken(token)

	//这里请求出来有两步, 如果tokenID==guestID, 直接返回, 默认全部关注
	//如果tokenID!=guestID, 还需要判断请求出来的user和token的关注关系
	userIDList, err := gorm.FindMultiConcern(guestID)
	if err != nil {
		log.Panicln("查找关注列表信息时失败")
		return resList, err
	}
	//获取到guestID关注的up的用户信息
	userInfoList, err := gorm.FindMultiUserInfo(userIDList)
	if err != nil {
		log.Panicln("查找关注用户的个人信息时失败")
		return resList, err
	}

	resList = make([]common.UserInfoResp, 0, len(userInfoList))
	for _, upInfo := range userInfoList {
		isFollowing := true
		if hostID != guestID {
			isFollowing = gorm.CheckFollowing(hostID, upInfo.ID)
		}
		resList = append(resList, util2.PackUserInfo(upInfo, isFollowing))
	}
	return resList, nil
}

// FollowerList 请求粉丝列表, 和请求关注列表逻辑一致, 不同的是需要更改一下查询数据库的信息
func FollowerList(token string, guestID uint) (resList []common.UserInfoResp, err error) {
	hostID := util2.GetHostIDFromToken(token)

	userIDList, err := gorm.FindMultiFollower(guestID)
	if err != nil {
		log.Panicln("查找粉丝列表信息时失败")
		return resList, err
	}
	//获取到guestID关注的up的用户信息
	userInfoList, err := gorm.FindMultiUserInfo(userIDList)
	if err != nil {
		log.Panicln("查找粉丝个人信息时失败")
		return resList, err
	}

	resList = make([]common.UserInfoResp, 0, len(userInfoList))
	for _, upInfo := range userInfoList {
		resList = append(resList, util2.PackUserInfo(upInfo, gorm.CheckFollowing(hostID, upInfo.ID)))
	}
	return resList, nil
}

// FriendList	这里要做一下过滤, 只有自己能看自己的好友列表, 当guestID与token不符时应当直接返回
func FriendList(token string, guestID uint) (resList []common.UserInfoResp, err error) {
	hostID := util2.GetHostIDFromToken(token)

	if hostID != guestID {
		return resList, errors.New("请求好友列表的id和token不一致, 不允许偷看别人的好友列表")
	}
	//分别获取粉丝和关注的up列表, 对其过滤, 同时存在的为好友
	fans, err := gorm.FindMultiFollower(guestID)
	if err != nil {
		log.Panicln("查找粉丝列表信息时失败")
		return resList, err
	}
	ups, err := gorm.FindMultiConcern(guestID)
	if err != nil {
		log.Panicln("查找关注列表信息时失败")
		return resList, err
	}

	mp := make(map[uint]interface{})
	for i := 0; i < len(fans); i++ {
		mp[fans[i]] = struct{}{}
	}
	friends := make([]uint, len(fans)>>1)
	for i := 0; i < len(ups); i++ {
		if _, ok := mp[ups[i]]; ok {
			friends = append(friends, ups[i])
		}
	}
	multiUserInfo, err := gorm.FindMultiUserInfo(friends)
	if err != nil {
		log.Panicln("获取好友信息时失败了")
		return resList, err
	}
	resList = make([]common.UserInfoResp, 0, len(multiUserInfo))
	for _, userInfo := range multiUserInfo {
		resList = append(resList, util2.PackUserInfo(userInfo, true))
	}
	return resList, nil
}
