package mapper

import "testing"

func TestGetMultiUserInfoCache(t *testing.T) {
	InitDBConnectorSupportTest()
	InitRedisConnector()
	GetMultiUserCache([]uint{1, 2})
}
