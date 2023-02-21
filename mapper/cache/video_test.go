package cache

import (
	"fmt"
	"testing"
	"tiktok/mapper"
	"time"
)

func TestSetInfoForRedis(t *testing.T) {
	_ = mapper.InitDBConnectorSupportTest()
	_ = mapper.InitRedisConnector()
	fmt.Println(time.Now().Unix())
}
