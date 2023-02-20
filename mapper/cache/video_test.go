package cache

import (
	"fmt"
	"testing"
	"tiktok/mapper"
	"time"
)

func TestSetInfoForRedis(t *testing.T) {
	mapper.InitDBConnectorSupportTest()
	mapper.InitRedisConnector()
	fmt.Println(time.Now().Unix())
}
