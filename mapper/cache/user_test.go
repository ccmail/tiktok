package cache

import (
	"testing"
	"tiktok/mapper"
)

func TestGetMultiUserInfoCache(t *testing.T) {
	mapper.InitDBConnectorSupportTest()
	mapper.InitRedisConnector()
	GetMultiUser(&[]uint{1, 2})
}
