package mapper

import (
	"github.com/redis/go-redis/v9"
	"log"
	"tiktok/config"
	"tiktok/pkg/constants"
	"tiktok/pkg/util"
)

func CheckFollowingCache(hostID, guestID uint) (ans bool, ok bool) {
	if hostID == 0 {
		return false, true
	}
	k := util.SpliceKey(constants.Follow, hostID, guestID)
	result, err := RedisConn.Get(RCtx, k).Result()
	if err != nil {
		log.Printf("缓存中不存在%v和%v关注关系\n", hostID, guestID)
		return false, false
	}
	return result == constants.RedisTrue, true
}

// CheckMultiFollowingCache 不存在关注信息的需要去mysql中查询
func CheckMultiFollowingCache(hostID uint, guestIDs *[]uint) (ans []bool, bad map[uint][]int) {
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

	result, err := RedisConn.MGet(RCtx, key...).Result()
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

func SetFollowingCache(hostID, guestID uint, isConcern bool) {
	k := util.SpliceKey(constants.Follow, hostID, guestID)
	val := constants.RedisTrue
	if !isConcern {
		val = constants.RedisFalse
	}
	err := RedisConn.Set(RCtx, k, val, config.RedisTimeout).Err()
	if err != nil {
		log.Println("插入缓存时失败")
	}
}

func SetMultiFollowingCache(hostID uint, guestID *[]uint, isFollow *[]bool) {
	if len(*guestID) != len(*isFollow) {
		log.Println("是否关注的数量和up的数量不一致, 插入失败!")
	}
	for i := 0; i < len(*guestID); i++ {
		SetFollowingCache(hostID, (*guestID)[i], (*isFollow)[i])
	}
	/*
		kv := make([]string, 0, len(*guestID)<<1)
		for i := 0; i < len(*guestID); i++ {
			v := constants.RedisTrue
			if !(*isFollow)[i] {
				v = constants.RedisFalse
			}
			kv = append(kv, util.SpliceKey(constants.Follow, hostID, (*guestID)[i]), v)
		}
		err := RedisConn.MSet(RCtx, kv).Err()
		if err != nil {
			log.Println("将关注信息写入缓存时失败...")
		}*/
}
