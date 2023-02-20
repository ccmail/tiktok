package redis

import (
	"testing"
	"tiktok/config"
)

func TestGetMultiUserInfoCache(t *testing.T) {
	config.InitDBConnectorSupportTest()
	config.InitRedisConnector()
	GetMultiUserCache([]uint{1, 2})
}
