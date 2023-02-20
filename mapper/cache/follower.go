package cache

import (
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"log"
	"tiktok/config"
	"tiktok/mapper"
	"tiktok/pkg/constants"
	"tiktok/pkg/util"
)

// CheckFollowing  判断 hostID 是否关注 guestID
func CheckFollowing(hostID, guestID uint) (ans bool, ok bool) {
	if hostID == 0 {
		return false, true
	}
	k := util.SpliceKey(constants.Follow, hostID, guestID)
	result, err := mapper.RedisConn.Get(RCtx, k).Result()
	if err != nil {
		log.Printf("缓存中不存在%v和%v关注关系\n", hostID, guestID)
		return false, false
	}
	return result == constants.RedisTrue, true
}

// CheckMultiFollowing 不存在关注信息的需要去mysql中查询
func CheckMultiFollowing(hostID uint, guestIDs *[]uint) (ans []bool, bad map[uint][]int) {
	ans = make([]bool, len(*guestIDs))
	//bad = make([]uint, len(*guestIDs))
	bad = make(map[uint][]int)
	if hostID == 0 {
		return ans, bad
	}

	key := make([]string, len(*guestIDs))
	for i := 0; i < len(key); i++ {
		key[i] = util.SpliceKey(constants.Follow, hostID, (*guestIDs)[i])
	}

	result, err := mapper.RedisConn.MGet(RCtx, key...).Result()
	if err != nil {
		log.Println("cache查询关注关系时出错")
		return ans, bad
	}

	for i := 0; i < len(ans); i++ {
		if result[i] == redis.Nil || result[i] == nil {
			if len(bad[(*guestIDs)[i]]) == 0 {
				bad[(*guestIDs)[i]] = make([]int, 0, 1)
			}
			bad[(*guestIDs)[i]] = append(bad[(*guestIDs)[i]], i)
			continue
		}
		ans[i] = result[i] == constants.RedisTrue
	}
	log.Printf("从cache中取出了%v个用户信息\n", len(ans)-len(bad))
	return ans, bad
}

func SetFollowing(hostID, guestID uint, isConcern bool) {
	if hostID == guestID {
		return
	}
	k := util.SpliceKey(constants.Follow, hostID, guestID)
	val := constants.RedisTrue
	if !isConcern {
		val = constants.RedisFalse
	}
	err := mapper.RedisConn.Set(RCtx, k, val, config.RedisTimeout).Err()
	if err != nil {
		log.Println("插入缓存时失败")
	}
}

func SetMultiFollowing(hostID uint, guestID *[]uint, isFollow *[]bool) {
	if len(*guestID) != len(*isFollow) {
		log.Println("是否关注的数量和up的数量不一致, 插入失败!")
	}
	for i := 0; i < len(*guestID); i++ {
		SetFollowing(hostID, (*guestID)[i], (*isFollow)[i])
	}
}

// GetFollowIDList 查询hostID的关注列表的id
func GetFollowIDList(hostID uint) []uint {
	key := util.SpliceKey(constants.Follow, hostID)
	result, err := mapper.RedisConn.Get(RCtx, key).Result()
	if err != nil {
		log.Printf("cache中不存在%v的关注信息", hostID)
		return nil
	}
	var ans []uint
	err = json.Unmarshal([]byte(result), &ans)
	return ans
}
func SetFollowIDList(hostID uint, followList *[]uint) {
	key := util.SpliceKey(constants.Follow, hostID)
	marshal, err := json.Marshal(*followList)
	if err != nil {
		log.Println("关注列表序列化失败, 无法插入法到cache")
	}
	mapper.RedisConn.Set(RCtx, key, marshal, config.RedisTimeout)
}

func GetFansIDList(hostID uint) []uint {
	key := util.SpliceKey(constants.Fans, hostID)
	result, err := mapper.RedisConn.Get(RCtx, key).Result()
	if err != nil {
		log.Printf("cache中不存在%v的关注信息", hostID)
		return nil
	}
	var ans []uint
	err = json.Unmarshal([]byte(result), &ans)
	return ans
}
func SetFansIDList(hostID uint, followList *[]uint) {
	key := util.SpliceKey(constants.Fans, hostID)
	marshal, err := json.Marshal(*followList)
	if err != nil {
		log.Println("粉丝列表序列化失败, 无法插入法到cache")
	}
	mapper.RedisConn.Set(RCtx, key, marshal, config.RedisTimeout)
}
