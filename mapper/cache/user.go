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
)

// GetMultiUser 返回值为在cache中查到的user和失败的user的ID对应在数组中的位置
func GetMultiUser(guestID *[]uint) (ans []model.User, failedUID map[uint][]int) {
	ans = make([]model.User, len(*guestID))
	failedUID = make(map[uint][]int)
	key := make([]string, len(*guestID))
	for i := range key {
		key[i] = util.SpliceKey(constants.User, (*guestID)[i])
	}
	result, err := mapper.RedisConn.MGet(RCtx, key...).Result()
	if err != nil {
		log.Println("cache查询用户信息时失败")
		return ans, failedUID
	}
	for i := 0; i < len(ans); i++ {
		if result[i] == redis.Nil || result[i] == nil {
			if len(failedUID[(*guestID)[i]]) == 0 {
				failedUID[(*guestID)[i]] = make([]int, 0, 1)
			}
			failedUID[(*guestID)[i]] = append(failedUID[(*guestID)[i]], i)
			continue
		}
		//log.Printf("%T,%v ", result[i], result[i])
		err = json.Unmarshal([]byte(result[i].(string)), &ans[i])
	}
	log.Printf("从cache中取出了%v个用户信息\n", len(ans)-len(failedUID))
	return ans, failedUID
}

func GetUser(id uint) (ans model.User, ok bool) {
	if id == 0 {
		return
	}
	k := util.SpliceKey(constants.User, id)
	result, err := mapper.RedisConn.Get(RCtx, k).Result()
	if err != nil {
		log.Printf("缓存中不存在%v的用户信息\n", id)
		return
	}
	err = json.Unmarshal([]byte(result), &ans)
	if err != nil {
		log.Println("解析缓存中存入的用户信息时出错")
		return
	}
	return ans, true
}

func SetMultiUser(users *[]model.User) {
	set := map[uint]model.User{}
	for i := 0; i < len(*users); i++ {
		if _, ok := set[(*users)[i].ID]; !ok {
			set[(*users)[i].ID] = (*users)[i]
		}
	}
	//MSet无法设置过期时间, Set和MSet对redis性能影响不大
	for _, user := range set {
		SetUser(&user)
	}
}

func SetUser(user *model.User) {
	key := util.SpliceKey(constants.User, user.ID)
	marshal, err := json.Marshal(&user)
	if err != nil {
		log.Println("用户信息序列化时失败, 失败原因为: ", err)
	}
	err = mapper.RedisConn.Set(RCtx, key, marshal, config.RedisTimeout).Err()
	if err != nil {
		log.Printf("用户%v写入缓存失败!失败原因为: %v\n", user.ID, err)
	}
}

func DelUser(id ...uint) {
	key := make([]string, len(id))
	for i := 0; i < len(id); i++ {
		key[i] = util.SpliceKey(constants.User, id[i])
	}
	mapper.RedisConn.Del(RCtx, key...)
}
