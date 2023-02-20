package mapper

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
	"tiktok/config"
	"tiktok/model"
	"tiktok/pkg/constants"
	"tiktok/pkg/util"
	"time"
)

var RCtx = context.Background()

type RedisOption func(rdb *redis.Client)

func GetVideoCache(key ...any) (v []model.Video) {
	k := util.SpliceKey(constants.Videos, key)
	result, err := RedisConn.Get(RCtx, k).Result()
	if err == redis.Nil {
		log.Println("redis中不存在", k)
		return v
	} else if err != nil {
		log.Println("其他的redis方面出错")
		return v
	}
	err = json.Unmarshal([]byte(result), &v)
	if err != nil {
		log.Println("redis反序列化时出错")
		return v
	}
	return v
}

func SetVideoCache(videos model.Video, key ...any) {
	k := util.SpliceKey(constants.Videos, key)
	marshal, err := json.Marshal(videos)
	if err != nil {
		log.Println("序列化存入redis时出错")
	}
	err = RedisConn.Set(RCtx, k, marshal, config.RedisTimeout).Err()
	if err != nil {
		log.Println("存入redis时出错")
	}
}

func SetMultiFeedCache(video *[]model.Video) {
	for _, v := range *video {
		AddFeedCache(v)
	}
}

// AddFeedCache 需要将视频写入feed流中的sortedSet
func AddFeedCache(video model.Video) {

	marshal, err := json.Marshal(video)
	if err != nil {
		log.Println("序列化video失败")
		return
	}
	k := util.SpliceKey(constants.Feed)
	add := RedisConn.ZAdd(RCtx, k, redis.Z{
		Score:  feedScore(video.CreatedAt, video.ID),
		Member: marshal,
	})

	if add.Err() != nil {
		log.Printf("将视%v频插入feed流失败,失败原因为%v\n", video.ID, add.Err())
	}
}

func GetFeedCache(latestTime time.Time) (ans []model.Video) {
	op := redis.ZRangeBy{
		Min:    "-1",
		Max:    strconv.FormatFloat(float64(latestTime.Unix())*(config.ZSetScoreUp), 'E', 10, 64),
		Offset: 0,
		Count:  int64(config.MaxFeedVideoCount),
	}
	//获取latestTime这个时间段内的视频
	k := util.SpliceKey(constants.Feed)
	result, err := RedisConn.ZRangeByScore(RCtx, k, &op).Result()
	if err != nil {
		log.Println("在cache查询feed流出错")
		return ans
	}
	//查询之后, 根据该latestTime删除一些成员, 设置策略为latestTime的前多少小时过期
	RedisConn.ZRemRangeByScore(RCtx, k,
		strconv.Itoa(-1),
		strconv.FormatFloat(
			float64(
				latestTime.Add(config.ZSetMemberTimeout).Unix())*(config.ZSetScoreUp),
			'E', 10, 64),
	)

	ans = make([]model.Video, 0, config.MaxFeedVideoCount)
	idx := 0
	for i := len(result) - 1; i >= 0; i-- {
		if idx >= config.MaxFeedVideoCount {
			break
		}
		temp := model.Video{}
		err = json.Unmarshal([]byte(result[i]), &temp)
		ans = append(ans, temp)
		idx++
	}
	log.Printf("从换种取到了%v个feed视频", len(ans))
	return ans
}

//time.Unix()
//使用videoID+time拼接score, 避免重复
//长度为10, 对应二进制应当为2^34,videoID这里进35位, 与之拼接
func feedScore(time time.Time, videoID uint) (x float64) {
	p := float64(videoID)
	x = float64(time.Unix())*config.ZSetScoreUp + (p / config.ZSetScoreDown)
	return x
}
