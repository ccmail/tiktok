package config

import "time"

// Key 使用时转化为[]byte
const Key string = "byte dance 11111 return"
const TokenLiveTime = 168 * time.Hour

// MaxFeedVideoCount 每次获取feed流的数量, 默认15, 最多设置为30
const MaxFeedVideoCount int = 15

const IllegalChar = "/\\:*?<>| "

const RedisTimeout = 15 * time.Minute

// ZSetScoreUp  用于拼接ZSet的分数, 计算时time.Time*ZSetScore+videoID/ZSetScoreDown
const ZSetScoreUp = 1e30

// ZSetScoreDown 将video作为小数拼接, 该值不建议超过1e15, 因为会丢失精度
const ZSetScoreDown = 1e10

const ZSetMemberTimeout = -24 * time.Hour

const MQLikeAdd = "like_add"
const MQLikeUpdate = "like_update"

// MQMaxLen 消息队列最大消息数量, 超过该值会删除, 0表示不设置
const MQMaxLen int64 = 10000
