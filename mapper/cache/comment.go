package cache

import (
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"log"
	"tiktok/config"
	"tiktok/mapper"
	"tiktok/model"
	"tiktok/pkg/constants"
	"tiktok/pkg/util"
	"time"
)

func SetComment(comment *model.Comment) {
	k := util.SpliceKey(constants.Comment, comment.VideoID)
	marshal, err := json.Marshal(comment)
	if err != nil {
		log.Println("评论信息序列化时失败, 无法插入cache", err)
	}
	err = mapper.RedisConn.ZAdd(RCtx, k, redis.Z{
		Score:  timeScore(time.Now(), comment.ID),
		Member: marshal,
	}).Err()
	if err != nil {
		log.Println("评论信息写入cache失败", err)
	}
	mapper.RedisConn.Expire(RCtx, k, config.RedisTimeout)
}

func GetCommentList(videoID uint) (ans []model.Comment) {
	k := util.SpliceKey(constants.Comment, videoID)
	result, err := mapper.RedisConn.ZRevRange(RCtx, k, 0, -1).Result()
	if err != nil {
		log.Println("cache中没有该video的相关信息")
	}
	ans = make([]model.Comment, len(result))
	for i := 0; i < len(result); i++ {
		_ = json.Unmarshal([]byte(result[i]), &ans[i])
	}
	//重新设置过期时间
	mapper.RedisConn.Expire(RCtx, k, config.RedisTimeout)
	return ans
}

func SetMultiComment(comment *[]model.Comment) {
	for i := 0; i < len(*comment); i++ {
		SetComment(&(*comment)[i])
	}
}

func DelComment(videoID model.Comment) {
	key := util.SpliceKey(constants.Comment, videoID)
	marshal, err := json.Marshal(videoID)
	if err != nil {
		log.Println("序列化时失误")
		return
	}
	mapper.RedisConn.ZRem(RCtx, key, marshal)
	if err != nil {
		//删除ZSet成员失败时, 之间将该ZSet直接删除
		mapper.RedisConn.Del(RCtx, key)
	} else {
		//否则更新下过期时间
		mapper.RedisConn.Expire(RCtx, key, config.RedisTimeout)
	}
}
