package mapper

import (
	"github.com/redis/go-redis/v9"
	"log"
	"tiktok/config"
	"tiktok/model"
	"tiktok/pkg/constants"
	"tiktok/pkg/util"
)

// CheckMultiFavoriteCache
// ans长度固定
// 检查cache中存在的点赞关系, 存在点赞关系的直接进行判断, 不存在的话, 添加到failedVID中, 后续直接取mysql查询failedVID中的数据, 并根据map索引写回到ans中
func CheckMultiFavoriteCache(uid uint, vInfo *[]model.Video) (ans []bool, bad map[uint][]int) {
	ans = make([]bool, len(*vInfo))
	bad = make(map[uint][]int)
	if uid == 0 {
		return ans, bad
	}

	key := make([]string, len(*vInfo))
	for i := 0; i < len(key); i++ {
		key[i] = util.SpliceKey(constants.Favorite, uid, (*vInfo)[i].ID)
	}

	result, err := RedisConn.MGet(RCtx, key...).Result()
	if err != nil {
		log.Println("cache中没有查到相关点赞记录")
		return ans, bad
	}

	for i := 0; i < len(ans); i++ {
		if result[i] == redis.Nil || result[i] == nil {
			if len(bad[(*vInfo)[i].ID]) == 0 {
				bad[(*vInfo)[i].ID] = make([]int, 0, 1)
			}
			bad[(*vInfo)[i].ID] = append(bad[(*vInfo)[i].ID], i)
			continue
		}
		ans[i] = result[i] == constants.RedisTrue
	}
	log.Printf("从cache中取出了%v个用户信息\n", len(ans)-len(bad))
	return ans, bad
}

// CheckFavoriteCache 查询cache中是否存在点赞信息,ans表示点赞情况, ok表示查询情况
func CheckFavoriteCache(uid uint, vInfo *model.Video) (ans, ok bool) {
	key := util.SpliceKey(constants.Favorite, uid, (*vInfo).ID)
	result, err := RedisConn.Get(RCtx, key).Result()
	if err != nil {
		log.Println("cache中没有查到相关点赞记录")
		return ans, false
	}
	return result == constants.RedisTrue, true

}

func SetMultiFavoriteCache(hostID uint, vInfos *[]model.Video, isLikes *[]bool) {
	if len(*vInfos) != len(*isLikes) {
		log.Println("传入的video信息和点赞信息的长度不同! 不能插入到缓存中")
	}
	for i := 0; i < len(*vInfos); i++ {
		SetFavoriteCache(hostID, (*vInfos)[i].ID, (*isLikes)[i])
	}
}

func SetFavoriteCache(uid, vid uint, isLike bool) {
	k := util.SpliceKey(constants.Favorite, uid, vid)
	v := constants.RedisTrue
	if !isLike {
		v = constants.RedisFalse
	}
	err := RedisConn.Set(RCtx, k, v, config.RedisTimeout).Err()
	if err != nil {
		log.Println("点赞信息插入缓存失败")
	}
}
