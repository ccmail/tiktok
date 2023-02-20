package mapper

import (
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"log"
	"tiktok/config"
	"tiktok/model"
	"tiktok/pkg/constants"
	"tiktok/pkg/util"
)

func GetMultiUserCache(guestID []uint) (ans []model.User, failedUID map[uint][]int) {
	ans = make([]model.User, len(guestID))
	failedUID = make(map[uint][]int)
	key := make([]string, len(guestID))
	for i := range key {
		key[i] = util.SpliceKey(constants.User, guestID[i])
	}
	result, err := RedisConn.MGet(RCtx, key...).Result()
	if err != nil {
		log.Println("cache查询用户信息时失败")
		return ans, failedUID
	}
	for i := 0; i < len(ans); i++ {
		if result[i] == redis.Nil || result[i] == nil {
			if len(failedUID[guestID[i]]) == 0 {
				failedUID[guestID[i]] = make([]int, 0, 1)
			}
			failedUID[guestID[i]] = append(failedUID[guestID[i]], i)
			continue
		}
		//log.Printf("%T,%v ", result[i], result[i])
		err = json.Unmarshal([]byte(result[i].(string)), &ans[i])
	}
	log.Printf("从cache中取出了%v个用户信息\n", len(ans)-len(failedUID))
	return ans, failedUID
}

func SetMultiUserCache(users *[]model.User) {
	set := map[uint]model.User{}
	for i := 0; i < len(*users); i++ {
		if _, ok := set[(*users)[i].ID]; !ok {
			set[(*users)[i].ID] = (*users)[i]
		}
	}
	//MSet无法设置过期时间, Set和MSet对redis性能影响不大
	for _, user := range set {
		SetUserCache(&user)
	}
}

func SetUserCache(user *model.User) {
	key := util.SpliceKey(constants.User, user.ID)
	marshal, err := json.Marshal(&user)
	if err != nil {
		log.Println("用户信息序列化时失败, 失败原因为: ", err)
	}
	err = RedisConn.Set(RCtx, key, marshal, config.RedisTimeout).Err()
	if err != nil {
		log.Printf("用户%v写入缓存失败!失败原因为: %v\n", user.ID, err)
	}
}
