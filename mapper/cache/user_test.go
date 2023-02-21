package cache

import (
	"testing"
	"tiktok/mapper"
)

func TestGetMultiUserInfoCache(t *testing.T) {
	_ = mapper.InitDBConnectorSupportTest()
	_ = mapper.InitRedisConnector()
	GetMultiUser(&[]uint{1, 2})
}
