package config

import "time"

// Key 使用时转化为[]byte
const Key string = "byte dance 11111 return"
const TokenLiveTime time.Duration = 168 * time.Hour

// MaxFeedVideoCount 每次获取feed流的数量, 默认15, 最多设置为30
const MaxFeedVideoCount int = 15

const IllegalChar = "/\\:*?<>| "

const RedisTimeout time.Duration = 15 * time.Minute

// ZSetScoreUp  用于拼接ZSet的分数, 计算时time.Time*ZSetScore+videoID/ZSetScoreDown
const ZSetScoreUp = 1e30

// ZSetScoreDown 将video作为小数拼接, 该值不建议超过1e15, 因为会丢失精度
const ZSetScoreDown = 1e10

const ZSetMemberTimeout = -24 * time.Hour
