package service

import (
	"log"
	"tiktok/mapper"
	"tiktok/pkg/common"
	"tiktok/pkg/middleware"
)

// Follow
// 1.向follow表中添加host和guest的关注关系
// 2.更改User表中对应的粉丝/关注up数量
func Follow(token string, guestID uint, isConcern bool) error {
	parseToken, err := middleware.ParseTokenCJS(token)
	if err != nil {
		return err
	}
	record, exist := mapper.ExistFollowRecord(parseToken.UserId, guestID)
	//如果不存在的话, 向follow插入信息, 并更新user表的字段
	if !exist {
		err := mapper.CreatFollowRecord(parseToken.UserId, guestID, isConcern)
		if err != nil {
			log.Panicln("插入数据库时发生错误")
			return err
		}
		//从未关注过的话直接插, 更新完毕之后return掉
		err = mapper.UpdateUserFollowCount(parseToken.UserId, guestID, isConcern)
		if err != nil {
			log.Panicln("更新粉丝/关注数量时出错")
			return err
		}
		return nil
		//mapper.UpdateUserFollowCount(parseToken.UserId, guestID)
	}

	if record.IsFollow == isConcern {
		//关注关系没有变化
		return nil
	}
	//更新host的关注数, 并增加host关注的up数量, 增加up的粉丝数量
	err = mapper.UpdateUserFollowCount(parseToken.UserId, guestID, isConcern)
	if err != nil {
		log.Panicln("更新粉丝/关注数量时出错")
		return err
	}
	//向follow表中添加follow关系记录
	err = mapper.UpdateFollowRecord(parseToken.UserId, guestID, isConcern)
	if err != nil {
		log.Panicln("更新关注关系时失败")
		return err
	}
	return nil
}

// FollowList 用户token去请求guestID对应的关注列表
func FollowList(token string, guestID uint) (returnList []common.UserInfoResp, err error) {
	hostID := getHostIDFromToken(token)
	//这里请求出来有两步, 如果tokenID==guestID, 直接返回, 默认全部关注
	//如果tokenID!=guestID, 还需要判断请求出来的user和token的关注关系
	userIDList, err := mapper.FindMultiConcern(guestID)
	if err != nil {
		return returnList, err
	}
	//获取到guestID关注的up的用户信息
	userInfoList, err := mapper.FindMultiUserInfo(userIDList)
	if err != nil {
		return returnList, err
	}

	returnList = make([]common.UserInfoResp, 0, len(userInfoList))
	for _, upInfo := range userInfoList {
		if hostID != guestID {
			returnList = append(returnList, packUserInfo(upInfo, hostID))
		} else {
			returnList = append(returnList, common.UserInfoResp{
				AuthorId:      upInfo.ID,
				Name:          upInfo.Name,
				FollowCount:   upInfo.FollowCount,
				FollowerCount: upInfo.FollowerCount,
				IsFollow:      true,
			})
		}
	}
	return returnList, nil
}
